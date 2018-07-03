package debug

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

type Cfg struct {
	Debug bool
	DebugFile string `yaml:"debug_file"`
}

func Config() Cfg {
	var c Cfg
	yamlFile, err := ioutil.ReadFile("data/config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}
	return c
}
