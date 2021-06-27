package taskfile

import (
	"strings"
)

// Cmd is a task command
type Cmd struct {
	Cmd    string
	Silent bool
	Loop   struct {
		Range    []int  //iterate start-to-end numbers
		Folder   string //iterate file names in the folder
		File     string //iterate line by line in the file
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

	var loopStruct struct {
		Range    []int
		Folder   string
		File     string
		Emit     string
		Run      string
		Parallel int
	}
	if err := unmarshal(&loopStruct); err != nil {
		return err
	}
	c.Loop.Range = loopStruct.Range
	c.Loop.Folder = loopStruct.Folder
	c.Loop.File = loopStruct.File
	c.Loop.Run = loopStruct.Run
	c.Loop.Parallel = loopStruct.Parallel

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
