package config

import "gopkg.in/yaml.v3"

type GCMerge struct {
	Global  GlobalT            `yaml:"global"`
	Configs map[string]ConfigT `yaml:"configs"`
}

type GlobalT struct {
	Type       string       `yaml:"type"`
	RawConfig  string       `yaml:"rawConfig"`
	Conditions []ConditionT `yaml:"conditions"`
	Actions    []ActionT    `yaml:"actions"`
}

type ConfigT struct {
	TargetConfig string       `yaml:"targetConfig"`
	RawConfig    string       `yaml:"rawConfig"`
	Conditions   []ConditionT `yaml:"conditions"`
	Actions      []ActionT    `yaml:"actions"`
}

type ConditionT struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
	Value    string `yaml:"value"`
}

type ActionT struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
	Script  string `yaml:"script"`
}

func Parse(config []byte) (cfg GCMerge, err error) {
	err = yaml.Unmarshal(config, &cfg)
	return cfg, err
}
