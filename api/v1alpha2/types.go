package v1alpha2

type CombiConfigT struct {
	Kind   string  `yaml:"kind"`
	Global GlobalT `yaml:"global"`
	Config ConfigT `yaml:"config"`
}

type GlobalT struct {
	Source     SourceT     `yaml:"source"`
	Conditions ConditionsT `yaml:"conditions"`
	Actions    ActionsT    `yaml:"actions"`
}

type ConfigT struct {
	MergedConfig string      `yaml:"mergedConfig"`
	Source       SourceT     `yaml:"source"`
	Conditions   ConditionsT `yaml:"conditions"`
	Actions      ActionsT    `yaml:"actions"`
}

type SourceT struct {
	Type       string      `yaml:"type,omitempty"`
	RawConfig  string      `yaml:"rawConfig,omitempty"`
	Filepath   string      `yaml:"filepath,omitempty"`
	Git        GitT        `yaml:"git,omitempty"`
	Kubernetes KubernetesT `yaml:"kubernetes,omitempty"`
}

type GitT struct {
	SshUrl         string `yaml:"sshUrl"`
	SshKeyFilepath string `yaml:"sshKeyFilepath"`
	Branch         string `yaml:"branch"`
	Filepath       string `yaml:"filepath"`
}

type KubernetesT struct {
	Kind      string `yaml:"kind"`
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	Key       string `yaml:"key"`
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
