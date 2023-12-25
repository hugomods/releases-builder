package config

type Language struct {
	Code string `yaml:"code"`
}

type Config struct {
	ContentDir   string     `yaml:"contentDir"`
	Languages    []Language `yaml:"languages"`
	Repositories []string   `yaml:"repositories"`
}

func NewConfig() *Config {
	return &Config{
		ContentDir: "content/releases",
	}
}
