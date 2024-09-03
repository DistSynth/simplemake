package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/gosuri/uitable"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	empty = struct{}{}
)

// RunTasks run all passed tasks
func RunTasks(cCtx *cli.Context) error {
	sm, err := readMakeFile()
	if err != nil {
		return err
	}

	taskMap := map[string]struct{}{}

	var execTask func(taskName string) error

	execTask = func(taskName string) error {
		log.Debug().Msgf("Executing task [%s]", taskName)
		task := sm.Tasks[taskName]

		for _, dep := range task.Dependencies {
			if _, ok := taskMap[dep]; ok {
				return fmt.Errorf("Loop detected")
			}
			taskMap[taskName] = empty
			if err := execTask(dep); err != nil {
				return err
			}
		}

		exists := true
		if len(task.Generates) == 0 {
			exists = false
		} else {
			log.Debug().Msg("Check task artifacts")
			for _, gen := range task.Generates {
				if _, err := os.Stat(gen); errors.Is(err, os.ErrNotExist) {
					log.Debug().Msgf("Task artifact [%s] doesn't exists.", gen)
					exists = false
					break
				} else {
					log.Debug().Msgf("Task artifact [%s] exists", gen)
				}
			}
		}

		if exists && !cCtx.Bool("force") {
			log.Info().Msgf("Skipping task [%s]", taskName)
			return nil
		}

		for k, v := range task.Var {
			log.Debug().Msgf("Task variable [%s=%s]", k, v)
		}

		for _, cmd := range task.Commands {
			log.Debug().Msgf("Executing task command [%s]", cmd)

			err = execute(cmd, task.Var)
			if err != nil {
				return fmt.Errorf("Task execution error [%s]", err.Error())
			}
		}

		return nil
	}

	for _, taskName := range cCtx.Args().Slice() {
		if _, ok := sm.Tasks[taskName]; ok {
			if err := execTask(taskName); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Task [%s] is not defined", taskName)
		}
	}

	return nil
}

// ListTasks list all tasks
func ListTasks(cCtx *cli.Context) error {
	log.Debug().Msg("Listing tasks")
	sm, err := readMakeFile()
	if err != nil {
		return err
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.Wrap = true
	table.AddRow("NAME", "DESCRIPTION")

	for name, task := range sm.Tasks {
		table.AddRow(name, task.Description)
	}

	fmt.Println(table)

	return nil
}

// ValidateMakeFile validate passed makefile
func ValidateMakeFile(cCtx *cli.Context) error {
	sm, err := readMakeFile()
	if err != nil {
		return err
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("NAME", "DESCRIPTION")

	for name, task := range sm.Tasks {
		table.AddRow(name, task.Description)
	}

	fmt.Println(table)

	return nil
}
