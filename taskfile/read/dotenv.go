package read

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/planetrio/planetr-compose/internal/compiler"
	"github.com/planetrio/planetr-compose/internal/templater"
	"github.com/planetrio/planetr-compose/taskfile"
)

func Dotenv(c compiler.Compiler, tf *taskfile.Taskfile, dir string) (*taskfile.Vars, error) {
	vars, err := c.GetTaskfileVariables()
	if err != nil {
		return nil, err
	}

	env := &taskfile.Vars{}

	tr := templater.Templater{Vars: vars, RemoveNoValue: true}

	for _, dotEnvPath := range tf.Dotenv {
		dotEnvPath = tr.Replace(dotEnvPath)

		if !filepath.IsAbs(dotEnvPath) {
			dotEnvPath = filepath.Join(dir, dotEnvPath)
		}
		if _, err := os.Stat(dotEnvPath); os.IsNotExist(err) {
			continue
		}

		envs, err := godotenv.Read(dotEnvPath)
		if err != nil {
			return nil, err
		}
		for key, value := range envs {
			if _, ok := env.Mapping[key]; !ok {
				env.Set(key, taskfile.Var{Static: value})
			}
		}
	}

	return env, nil
}
