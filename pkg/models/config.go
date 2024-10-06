package models

type CompanyConfig struct {
	ID      int64    `yaml:"id"`
	Name    string   `yaml:"name"`
	Columns []string `yaml:"columns"`
}

type Config struct {
	Companies []CompanyConfig `yaml:"companies"`
}
