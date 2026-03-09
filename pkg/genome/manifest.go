package genome

type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Manifest struct {
	Project Project           `yaml:"project"`
	Traits  []string          `yaml:"traits"`
	Modules map[string]string `yaml:"modules"`
}
