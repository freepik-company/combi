package raw

import (
	"combi/api/v1alpha2"
	"os"
	"reflect"
)

type RawSourceT struct {
	RawConfig    string
	StoredConfig []byte
}

func (s *RawSourceT) Init(source v1alpha2.SourceT) (err error) {
	s.RawConfig = source.RawConfig
	return err
}

func (s *RawSourceT) GetConfig() (config []byte, updated bool, err error) {
	config = []byte(os.ExpandEnv(s.RawConfig))

	if !reflect.DeepEqual(s.StoredConfig, config) {
		updated = true
	}

	return config, updated, err
}
