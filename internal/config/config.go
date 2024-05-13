package config

import "gopkg.in/yaml.v3"

type Gcu struct {
	Global  GlobalT            `yaml:"global"`
	Configs map[string]ConfigT `yaml:"configs"`
}

type GlobalT struct {
	Raw        string            `yaml:"raw"`
	Fields     map[string]string `yaml:"fields"`
	Conditions []ConditionT      `yaml:"conditions"`
	Actions    []ActionT         `yaml:"actions"`
}

type ConfigT struct {
	Raw        string            `yaml:"raw"`
	Fields     map[string]string `yaml:"fields"`
	Conditions []ConditionT      `yaml:"conditions"`
	Actions    []ActionT         `yaml:"actions"`
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

func Parse(config []byte) (cfg Gcu, err error) {
	err = yaml.Unmarshal(config, &cfg)
	return cfg, err
}
