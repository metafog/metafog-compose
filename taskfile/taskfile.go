package taskfile

// Taskfile represents a Taskfile.yml
type Taskfile struct {
	Version    string
	Expansions int
	Output     string
	Method     string
	Includes   *IncludedTaskfiles
	Vars       *Vars
	Env        *Vars
	Tasks      Tasks
	Silent     bool
	Dotenv     []string
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (tf *Taskfile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var taskfile struct {
		Version    string
		Expansions int
		Output     string
		Method     string
		Includes   *IncludedTaskfiles
		Vars       *Vars
		Env        *Vars
		Tasks      Tasks
		Silent     bool
		Dotenv     []string
	}
	if err := unmarshal(&taskfile); err != nil {
		return err
	}
	tf.Version = taskfile.Version
	tf.Expansions = taskfile.Expansions
	tf.Output = taskfile.Output
	tf.Method = taskfile.Method
	tf.Includes = taskfile.Includes
	tf.Vars = taskfile.Vars
	tf.Env = taskfile.Env
	tf.Tasks = taskfile.Tasks
	tf.Silent = taskfile.Silent
	tf.Dotenv = taskfile.Dotenv
	if tf.Expansions <= 0 {
		tf.Expansions = 2
	}
	if tf.Vars == nil {
		tf.Vars = &Vars{}
	}
	if tf.Env == nil {
		tf.Env = &Vars{}
	}
	return nil
}
