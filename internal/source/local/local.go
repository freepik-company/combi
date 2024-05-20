package local

import (
	"combi/internal/flags"
	"os"
	"reflect"
)

type LocalT struct {
	ConfigFilepath string
	StoredConfig   []byte
	IsUpdated      bool
}

func (s *LocalT) Init(f flags.DaemonFlagsT) {
	s.ConfigFilepath = f.SourcePath
}

func (s *LocalT) GetConfig() (config []byte, err error) {
	s.IsUpdated = false

	if config, err = os.ReadFile(s.ConfigFilepath); err != nil {
		return config, err
	}

	if !reflect.DeepEqual(s.StoredConfig, config) {
		s.IsUpdated = true
		s.StoredConfig = config
	}

	return config, err
}

func (s *LocalT) NeedUpdate() bool {
	return s.IsUpdated
}
