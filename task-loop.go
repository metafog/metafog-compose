package task

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ahmetask/worker"
	"github.com/go-stomp/stomp"
	"github.com/go-task/task/v3/taskfile"
	"github.com/radovskyb/watcher"
	kafka "github.com/segmentio/kafka-go"
)

const (
	LOOP_ARG          = "ARG"
	LOOP_INDX         = "INDX"
	KAFKA_GROUP       = "planetr-group"
	JOB_QUEUE_SIZE    = 10
	ACTIVEMQ_PROTOCOL = "tcp"
)

type ParallelTask struct {
	Ctx  context.Context
	Call taskfile.Call
	e    *Executor
}

/*implement work interface*/
func (j *ParallelTask) Do() {
	err := j.e.RunTask(j.Ctx, j.Call)
	if err != nil {
		fmt.Println("Err: ", err)
	}
}

func (e *Executor) loopTasks(ctx context.Context, cmd *taskfile.Cmd) error {

	threads := cmd.Loop.Parallel
	if threads < 1 {
		threads = 1
	}

	// Initialize Pool
	pool := worker.NewWorkerPool(threads, JOB_QUEUE_SIZE)
	pool.Start()

	//Stop Worker Pool
	defer pool.Stop()

	switch {
	case len(cmd.Loop.Range) == 2:
		return e.runLoopRange(ctx, cmd, pool)
	case len(cmd.Loop.Folder) > 0:
		return e.runLoopFolder(ctx, cmd, pool)
	case len(cmd.Loop.FolderWatch) > 0:
		return e.runLoopFolderWatch(ctx, cmd, pool)
	case len(cmd.Loop.File) > 0:
		return e.runLoopFile(ctx, cmd, pool)
	case cmd.Loop.Timer > 0:
		return e.runLoopTimer(ctx, cmd, pool)
	case len(cmd.Loop.Activemq) == 2:
		return e.runLoopActiveMQ(ctx, cmd, pool)
	case len(cmd.Loop.Kafka) == 2:
		return e.runLoopKafka(ctx, cmd, pool)
	}

	return nil
}

func (e *Executor) runLoopRange(ctx context.Context, cmd *taskfile.Cmd, pool *worker.Pool) error {
	fmt.Println("Loop > Range", cmd.Loop.Range[0], cmd.Loop.Range[1])

	indx := 0
	for i := cmd.Loop.Range[0]; i <= cmd.Loop.Range[1]; i++ {
		//ARGS
		vars := e.addVars(ctx, cmd, strconv.Itoa(i), indx)
		ptask := ParallelTask{
			Ctx:  ctx,
			Call: taskfile.Call{Task: cmd.Loop.Run, Vars: vars},
			e:    e,
		}
		pool.Submit(&ptask)
		indx++
	}
	return nil
}

func (e *Executor) runLoopFolder(ctx context.Context, cmd *taskfile.Cmd, pool *worker.Pool) error {
	fmt.Println("Loop > Folder", cmd.Loop.Folder)

	items, _ := ioutil.ReadDir(cmd.Loop.Folder)
	indx := 0
	for _, item := range items {
		if !item.IsDir() {
			//ARGS
			vars := e.addVars(ctx, cmd, item.Name(), indx)
			ptask := ParallelTask{
				Ctx:  ctx,
				Call: taskfile.Call{Task: cmd.Loop.Run, Vars: vars},
				e:    e,
			}
			pool.Submit(&ptask)
			indx++
		}
	}
	return nil
}

func (e *Executor) runLoopFolderWatch(ctx context.Context, cmd *taskfile.Cmd, pool *worker.Pool) error {
	fmt.Println("Loop > FolderWatch", cmd.Loop.FolderWatch, "...")

	w := watcher.New()
	w.SetMaxEvents(1)

	// Only notify create events.
	w.FilterOps(watcher.Create)

	go func() {
		indx := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Interrupt")
				os.Exit(1)
				return
			case event := <-w.Event:
				//ARGS
				vars := e.addVars(ctx, cmd, event.Path, indx)
				ptask := ParallelTask{
					Ctx:  ctx,
					Call: taskfile.Call{Task: cmd.Loop.Run, Vars: vars},
					e:    e,
				}
				pool.Submit(&ptask)
				indx++

			case err := <-w.Error:
				fmt.Println("ERROR", err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.AddRecursive(cmd.Loop.FolderWatch); err != nil {
		return err
	}

	// Start the watching process - it'll check for changes every 1s.
	if err := w.Start(time.Millisecond * 1000); err != nil {
		return err
	}

	return nil
}

func (e *Executor) runLoopFile(ctx context.Context, cmd *taskfile.Cmd, pool *worker.Pool) error {
	fmt.Println("Loop > File", cmd.Loop.File)

	file, err := os.Open(cmd.Loop.File)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	indx := 0
	for scanner.Scan() {
		//ARGS
		vars := e.addVars(ctx, cmd, scanner.Text(), indx)
		ptask := ParallelTask{
			Ctx:  ctx,
			Call: taskfile.Call{Task: cmd.Loop.Run, Vars: vars},
			e:    e,
		}
		pool.Submit(&ptask)
		indx++
	}
	return nil
}

func (e *Executor) runLoopTimer(ctx context.Context, cmd *taskfile.Cmd, pool *worker.Pool) error {
	fmt.Println("Loop > Timer", time.Second, "...")

	ticker := time.NewTicker(time.Second * time.Duration(cmd.Loop.Timer))
	defer ticker.Stop()

	indx := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Interrupt")
			return nil
		case <-ticker.C:
			//ARGS
			arg := time.Now().String()
			vars := e.addVars(ctx, cmd, arg, indx)
			ptask := ParallelTask{
				Ctx:  ctx,
				Call: taskfile.Call{Task: cmd.Loop.Run, Vars: vars},
				e:    e,
			}
			pool.Submit(&ptask)
			indx++
		}
	}
}

func (e *Executor) runLoopActiveMQ(ctx context.Context, cmd *taskfile.Cmd, pool *worker.Pool) error {
	fmt.Println("Loop > ActiveMQ", cmd.Loop.Activemq[0], "...")

	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.HeartBeat(7200*time.Second, 7200*time.Second),
		stomp.ConnOpt.HeartBeatError(360 * time.Second),
	}

	conn, err := stomp.Dial(ACTIVEMQ_PROTOCOL, cmd.Loop.Activemq[0], options...)
	if err != nil {
		fmt.Println("Cannot connect to server", err.Error())
		return err
	}

	sub, err := conn.Subscribe(cmd.Loop.Activemq[1], stomp.AckAuto)
	if err != nil {
		fmt.Println("Cannot subscribe to topic", err.Error())
		return err
	}

	indx := 0
	for {
		msg := <-sub.C
		if msg != nil {
			vars := e.addVars(ctx, cmd, string(msg.Body), indx)
			ptask := ParallelTask{
				Ctx:  ctx,
				Call: taskfile.Call{Task: cmd.Loop.Run, Vars: vars},
				e:    e,
			}
			pool.Submit(&ptask)
			indx++
		}
	}
}

func (e *Executor) runLoopKafka(ctx context.Context, cmd *taskfile.Cmd, pool *worker.Pool) error {
	fmt.Println("Loop > Kafka", cmd.Loop.Kafka[0], cmd.Loop.Kafka[1], "...")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cmd.Loop.Kafka[0]},
		Topic:   cmd.Loop.Kafka[1],
		GroupID: KAFKA_GROUP,
	})

	indx := 0
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Println("could not read message " + err.Error())
		}

		vars := e.addVars(ctx, cmd, string(msg.Value), indx)
		ptask := ParallelTask{
			Ctx:  ctx,
			Call: taskfile.Call{Task: cmd.Loop.Run, Vars: vars},
			e:    e,
		}
		pool.Submit(&ptask)
		indx++
	}
}

func (e *Executor) addVars(ctx context.Context, cmd *taskfile.Cmd, arg string, indx int) *taskfile.Vars {
	vars := &taskfile.Vars{}
	vars.Mapping = make(map[string]taskfile.Var)

	vars.Keys = append(vars.Keys, LOOP_ARG)
	vars.Mapping[LOOP_ARG] = taskfile.Var{Static: arg}

	vars.Keys = append(vars.Keys, LOOP_INDX)
	vars.Mapping[LOOP_INDX] = taskfile.Var{Static: strconv.Itoa(indx + 1)}

	return vars
}
