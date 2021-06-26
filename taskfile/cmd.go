package taskfile

import (
	"strings"
)

// Cmd is a task command
type Cmd struct {
	Cmd    string
	Silent bool
	Loop   struct {
		Range    []int
		Folder   string
		File     string
		Run      string
		Parallel int
	}
	Dcurun      string
	Task        string
	Vars        *Vars
	IgnoreError bool
}

// Dep is a task dependency
type Dep struct {
	Task string
	Vars *Vars
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (c *Cmd) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var cmd string
	if err := unmarshal(&cmd); err == nil {
		if strings.HasPrefix(cmd, "^") {
			c.Task = strings.TrimPrefix(cmd, "^")
		} else {
			c.Cmd = cmd
		}
		return nil
	}

	var cmdStruct struct {
		Cmd         string
		Silent      bool
		IgnoreError bool `yaml:"ignore_error"`
	}
	if err := unmarshal(&cmdStruct); err == nil && cmdStruct.Cmd != "" {
		c.Cmd = cmdStruct.Cmd
		c.Silent = cmdStruct.Silent
		c.IgnoreError = cmdStruct.IgnoreError
		return nil
	}

	//LOOP BEG
	var looprangeStruct struct {
		Range    []int
		Folder   string
		File     string
		Emit     string
		Run      string
		Parallel int
	}
	if err := unmarshal(&looprangeStruct); err != nil {
		return err
	}
	c.Loop.Range = looprangeStruct.Range
	c.Loop.Folder = looprangeStruct.Folder
	c.Loop.File = looprangeStruct.File
	c.Loop.Run = looprangeStruct.Run
	c.Loop.Parallel = looprangeStruct.Parallel
	//LOOP END

	var taskCall struct {
		Task string
		Vars *Vars
	}
	if err := unmarshal(&taskCall); err != nil {
		return err
	}
	c.Task = taskCall.Task
	c.Vars = taskCall.Vars

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (d *Dep) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var task string
	if err := unmarshal(&task); err == nil {
		d.Task = task
		return nil
	}
	var taskCall struct {
		Task string
		Vars *Vars
	}
	if err := unmarshal(&taskCall); err != nil {
		return err
	}
	d.Task = taskCall.Task
	d.Vars = taskCall.Vars
	return nil
}
