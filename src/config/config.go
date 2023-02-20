package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Cache            bool   `yaml:"cache"`
		IntervalTimeType string `yaml:"interval_time_type"`
		IntervalAmount   int    `yaml:"interval_amount"`
	} `yaml:"server"`
}

func GetConfig(fileName string) Config {
	config := Config{}

	// Read YAML file
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal YAML data into config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config

}
