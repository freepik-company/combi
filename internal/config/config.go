package config

import (
	"combi/api/v1alpha2"

	"gopkg.in/yaml.v3"
)

func Parse(config []byte) (cfg v1alpha2.CombiConfigT, err error) {
	err = yaml.Unmarshal(config, &cfg)
	return cfg, err
}
