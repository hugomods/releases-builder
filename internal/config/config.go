package config

type Language struct {
	Code string `yaml:"code"`
}

type Params struct {
	Authors    []string `yaml:"authors"`
	Categories []string `yaml:"categories"`
	Images     []string `yaml:"images"`
	Series     []string `yaml:"series"`
	Tags       []string `yaml:"tags"`
}

type Config struct {
	ContentDir   string     `yaml:"contentDir"`
	Languages    []Language `yaml:"languages"`
	Repositories []string   `yaml:"repositories"`
	Params       Params     `yaml:"params"`
}

func NewConfig() *Config {
	return &Config{
		ContentDir: "content/releases",
	}
}
