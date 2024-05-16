package config

import (
	"gcmerge/api/v1alpha1"

	"gopkg.in/yaml.v3"
)

func Parse(config []byte) (cfg v1alpha1.GCMerge, err error) {
	err = yaml.Unmarshal(config, &cfg)
	return cfg, err
}
