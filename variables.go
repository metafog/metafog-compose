package task

import (
	"path/filepath"
	"strings"

	"github.com/planetrio/planetr-compose/internal/execext"
	"github.com/planetrio/planetr-compose/internal/status"
	"github.com/planetrio/planetr-compose/internal/templater"
	"github.com/planetrio/planetr-compose/taskfile"
)

// CompiledTask returns a copy of a task, but replacing variables in almost all
// properties using the Go template package.
func (e *Executor) CompiledTask(call taskfile.Call) (*taskfile.Task, error) {
	return e.compiledTask(call, true)
}

// FastCompiledTask is like CompiledTask, but it skippes dynamic variables.
func (e *Executor) FastCompiledTask(call taskfile.Call) (*taskfile.Task, error) {
	return e.compiledTask(call, false)
}

func (e *Executor) compiledTask(call taskfile.Call, evaluateShVars bool) (*taskfile.Task, error) {
	origTask, ok := e.Taskfile.Tasks[call.Task]
	if !ok {
		return nil, &taskNotFoundError{call.Task}
	}

	var vars *taskfile.Vars
	var err error
	if evaluateShVars {
		vars, err = e.Compiler.GetVariables(origTask, call)
	} else {
		vars, err = e.Compiler.FastGetVariables(origTask, call)
	}
	if err != nil {
		return nil, err
	}

	r := templater.Templater{Vars: vars, RemoveNoValue: true}

	new := taskfile.Task{
		Task:        origTask.Task,
		Label:       r.Replace(origTask.Label),
		Desc:        r.Replace(origTask.Desc),
		Summary:     r.Replace(origTask.Summary),
		Sources:     r.ReplaceSlice(origTask.Sources),
		Generates:   r.ReplaceSlice(origTask.Generates),
		Dir:         r.Replace(origTask.Dir),
		Vars:        call.Vars,
		Env:         nil,
		Silent:      origTask.Silent,
		Method:      r.Replace(origTask.Method),
		Prefix:      r.Replace(origTask.Prefix),
		IgnoreError: origTask.IgnoreError,
	}
	new.Dir, err = execext.Expand(new.Dir)
	if err != nil {
		return nil, err
	}
	if e.Dir != "" && !filepath.IsAbs(new.Dir) {
		new.Dir = filepath.Join(e.Dir, new.Dir)
	}
	if new.Prefix == "" {
		new.Prefix = new.Task
	}

	new.Env = &taskfile.Vars{}
	new.Env.Merge(r.ReplaceVars(e.Taskfile.Env))
	new.Env.Merge(r.ReplaceVars(origTask.Env))
	if evaluateShVars {
		err = new.Env.Range(func(k string, v taskfile.Var) error {
			static, err := e.Compiler.HandleDynamicVar(v, new.Dir)
			if err != nil {
				return err
			}
			new.Env.Set(k, taskfile.Var{Static: static})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if len(origTask.Cmds) > 0 {
		new.Cmds = make([]*taskfile.Cmd, len(origTask.Cmds))
		for i, cmd := range origTask.Cmds {
			new.Cmds[i] = &taskfile.Cmd{
				Task:        r.Replace(cmd.Task),
				Silent:      cmd.Silent,
				Cmd:         r.Replace(cmd.Cmd),
				Vars:        r.ReplaceVars(cmd.Vars),
				IgnoreError: cmd.IgnoreError,
			}

			new.Cmds[i].Loop.Range = cmd.Loop.Range
			new.Cmds[i].Loop.List = cmd.Loop.List
			new.Cmds[i].Loop.Folder = cmd.Loop.Folder
			new.Cmds[i].Loop.FolderWatch = cmd.Loop.FolderWatch
			new.Cmds[i].Loop.File = cmd.Loop.File
			new.Cmds[i].Loop.Timer = cmd.Loop.Timer
			new.Cmds[i].Loop.Activemq = cmd.Loop.Activemq
			new.Cmds[i].Loop.Kafka = cmd.Loop.Kafka
			new.Cmds[i].Loop.Run = cmd.Loop.Run
			new.Cmds[i].Loop.Parallel = cmd.Loop.Parallel
		}
	}
	if len(origTask.Deps) > 0 {
		new.Deps = make([]*taskfile.Dep, len(origTask.Deps))
		for i, dep := range origTask.Deps {
			new.Deps[i] = &taskfile.Dep{
				Task: r.Replace(dep.Task),
				Vars: r.ReplaceVars(dep.Vars),
			}
		}
	}

	if len(origTask.Preconditions) > 0 {
		new.Preconditions = make([]*taskfile.Precondition, len(origTask.Preconditions))
		for i, precond := range origTask.Preconditions {
			new.Preconditions[i] = &taskfile.Precondition{
				Sh:  r.Replace(precond.Sh),
				Msg: r.Replace(precond.Msg),
			}
		}
	}

	if len(origTask.Status) > 0 {
		for _, checker := range []status.Checker{e.timestampChecker(&new), e.checksumChecker(&new)} {
			value, err := checker.Value()
			if err != nil {
				return nil, err
			}
			vars.Set(strings.ToUpper(checker.Kind()), taskfile.Var{Live: value})
		}

		// Adding new variables, requires us to refresh the templaters
		// cache of the the values manually
		r.ResetCache()

		new.Status = r.ReplaceSlice(origTask.Status)
	}

	return &new, r.Err()
}
