// +build elasticsearch

package core

import (
	"github.com/olivere/elastic/v7"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

type ESConfig struct {
	Host                []string `yaml:"host"`
	Sniff               bool     `yaml:"sniff"`
	HealthcheckInterval int      `yaml:"healtcheck_interval"`
}

// Info from config file
type ProjectESConfig struct {
	ES *ESConfig `yaml:"es"`
}

func (ucc *ProjectESConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawStuff ProjectESConfig
	raw := rawStuff{
		ES: &ESConfig{
			Host:                []string{"http://127.0.0.1:9200"},
			Sniff:               false,
			HealthcheckInterval: 5,
		},
	}
	// Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*ucc = ProjectESConfig(raw)
	return nil

}

func NewProjectESClient() *elastic.Client {
	opts := make([]elastic.ClientOptionFunc, 0)
	config := &ProjectESConfig{}
	err := yaml.Unmarshal(CurrentConfig.ConfigContent, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	opts = append(opts, elastic.SetURL(config.ES.Host...))
	opts = append(opts, elastic.SetSniff(config.ES.Sniff))
	opts = append(opts, elastic.SetHealthcheckInterval(time.Duration(config.ES.HealthcheckInterval)*time.Second))
	opts = append(opts, elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)))
	opts = append(opts, elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	return client
}
