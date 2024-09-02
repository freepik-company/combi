package file

import (
	"combi/api/v1alpha2"
	"os"
	"reflect"
)

type FileSourceT struct {
	ConfigFilepath string
	StoredConfig   []byte
}

func (s *FileSourceT) Init(source v1alpha2.SourceT) {
	s.ConfigFilepath = source.Filepath
}

func (s *FileSourceT) GetConfig() (config []byte, updated bool, err error) {
	if config, err = os.ReadFile(s.ConfigFilepath); err != nil {
		return config, updated, err
	}
	config = []byte(os.ExpandEnv(string(config)))

	if !reflect.DeepEqual(s.StoredConfig, config) {
		updated = true
	}

	return config, updated, err
}
