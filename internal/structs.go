package internal

type SimpleMake struct {
	Tasks map[string]MakeTask `yaml:"tasks"`
}

type MakeTask struct {
	Description  string            `yaml:"description"`
	Dependencies []string          `yaml:"dependencies"`
	Commands     []string          `yaml:"commands"`
	Generates    []string          `yaml:"generates"`
	Var          map[string]string `yaml:"var"`
}
