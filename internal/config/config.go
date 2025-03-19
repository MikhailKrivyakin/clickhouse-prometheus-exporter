package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Servers []Server `yaml:"servers"`
	Queries []Query  `yaml:"queries"`
}

type Server struct {
	Name string `yaml:"name"`
	DSN  string `yaml:"dsn"`
}

type Query struct {
	Name        string   `yaml:"name"`
	Query       string   `yaml:"query"`
	MetricName  string   `yaml:"metric_name"`
	Help        string   `yaml:"help"`
	Type        string   `yaml:"type"`
	ValueColumn string   `yaml:"value_column"` // Столбец, используемый как значение метрики
	Labels      []string `yaml:"labels"`       // Столбцы, используемые как лейблы
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
