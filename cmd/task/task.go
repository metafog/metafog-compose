package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/pflag"

	task "github.com/planetrio/planetr-compose"
	"github.com/planetrio/planetr-compose/args"
	"github.com/planetrio/planetr-compose/internal/logger"
	"github.com/planetrio/planetr-compose/taskfile"
)

const usage = `Usage: planetr-compose [-t taskfile] [task...]

Runs the specified task(s). Falls back to the "default" task if no task name
was specified, or lists all tasks if an unknown task name was specified.

Example: 'planetr-compose hello' with the following 'Taskfile.yml' file will generate an
'output.txt' file with the content "hello".

Credits: https://github.com/go-task/task

'''
tasks:
  hello:
    cmds:
      - echo "I am going to write a file named 'output.txt' now."
      - echo "hello" > output.txt
    generates:
      - output.txt
'''

Options:
`

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	pflag.Usage = func() {
		log.Print(usage)
		pflag.PrintDefaults()
	}

	var (
		versionFlag bool
		helpFlag    bool
		list        bool
		silent      bool
		entrypoint  string
	)

	pflag.StringVarP(&entrypoint, "taskfile", "t", "", `choose which Taskfile to run. Defaults to "Taskfile.yml"`)
	pflag.BoolVarP(&versionFlag, "version", "v", false, "show version")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "shows usage")
	pflag.BoolVarP(&list, "list", "l", false, "lists tasks with description of current Taskfile")
	pflag.BoolVarP(&silent, "silent", "s", false, "disables echoing")
	pflag.Parse()

	if versionFlag {
		fmt.Printf("planetr-compose version: %s\n", getVersion())
		return
	}

	if helpFlag {
		pflag.Usage()
		return
	}

	if entrypoint != "" {
		entrypoint = filepath.Base(entrypoint)
	} else {
		entrypoint = "Taskfile.yml"
	}

	e := task.Executor{
		Verbose:    false,
		Silent:     silent,
		Entrypoint: entrypoint,
		Color:      true,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := e.Setup(); err != nil {
		log.Fatal(err)
	}

	if list {
		e.PrintTasksHelp()
		return
	}

	var (
		calls                 []taskfile.Call
		globals               *taskfile.Vars
		tasksAndVars, cliArgs = getArgs()
	)

	calls, globals = args.ParseV3(tasksAndVars...)

	globals.Set("CLI_ARGS", taskfile.Var{Static: strings.Join(cliArgs, " ")})
	e.Taskfile.Vars.Merge(globals)

	watch := false

	ctx := context.Background()
	if !watch {
		ctx = getSignalContext()
	}

	tBeg := time.Now()
	if err := e.Run(ctx, calls...); err != nil {
		e.Logger.Errf(logger.Red, "%v", err)
		os.Exit(1)
	}
	tEnd := time.Now()

	diff := tEnd.Sub(tBeg)
	fmt.Println("Finished in:", diff)
}

func getArgs() (tasksAndVars, cliArgs []string) {
	var (
		args          = pflag.Args()
		doubleDashPos = pflag.CommandLine.ArgsLenAtDash()
	)

	if doubleDashPos != -1 {
		tasksAndVars = args[:doubleDashPos]
		cliArgs = args[doubleDashPos:]
	} else {
		tasksAndVars = args
	}

	return
}

func getSignalContext() context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-ch
		log.Printf("task: signal received: %s", sig)
		cancel()
		os.Exit(1)
	}()
	return ctx
}

func getVersion() string {
	return "1.0"
}
