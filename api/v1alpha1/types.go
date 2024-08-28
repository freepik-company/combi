package v1alpha1

type GCMerge struct {
	Kind    string             `yaml:"kind"`
	Global  GlobalT            `yaml:"global"`
	Configs map[string]ConfigT `yaml:"configs"`
}

type GlobalT struct {
	RawConfig  string      `yaml:"rawConfig"`
	Conditions ConditionsT `yaml:"conditions"`
	Actions    ActionsT    `yaml:"actions"`
}

type ConfigT struct {
	TargetConfig string      `yaml:"targetConfig"`
	MergedConfig string      `yaml:"mergedConfig"`
	RawConfig    string      `yaml:"rawConfig"`
	Conditions   ConditionsT `yaml:"conditions"`
	Actions      ActionsT    `yaml:"actions"`
}

type ConditionsT struct {
	Mandatory []ConditionT `yaml:"mandatory"`
	Optional  []ConditionT `yaml:"optional"`
}

type ActionsT struct {
	OnSuccess []ActionT `yaml:"onSuccess"`
	OnFailure []ActionT `yaml:"onFailure"`
}

type ConditionT struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
	Value    string `yaml:"value"`
}

type ActionT struct {
	Name    string   `yaml:"name"`
	Command []string `yaml:"command"`
	Script  string   `yaml:"script"`
}
