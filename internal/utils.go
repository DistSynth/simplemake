package internal

import (
	"bytes"
	"os"
	"os/exec"
	"text/template"

	"github.com/flynn/go-shlex"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func readMakeFile() (*SimpleMake, error) {
	log.Debug().Msg("Reading makefile")
	yamlFile, err := os.ReadFile("simplemake.yaml")
	if err != nil {
		return nil, err
	}
	sm := SimpleMake{}
	err = yaml.Unmarshal(yamlFile, &sm)
	if err != nil {
		return nil, err
	}
	return &sm, nil
}

func execute(name string, env map[string]string) error {
	tmpl, err := template.New("test").Parse(name)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, env); err != nil {
		return err
	}

	result := buf.String()

	args, _ := shlex.Split(result)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
