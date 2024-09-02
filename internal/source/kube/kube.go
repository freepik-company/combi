package kube

import (
	"reflect"

	"combi/api/v1alpha2"
	"combi/internal/logger"
)

type KubeSourceT struct {
	StoredConfig []byte
}

func (s *KubeSourceT) Init(source v1alpha2.SourceT) {
	logger.Log.Fatalf("kube source type not implemented yet")
}

func (s *KubeSourceT) GetConfig() (config []byte, updated bool, err error) {

	if !reflect.DeepEqual(s.StoredConfig, config) {
		updated = true
	}

	return config, updated, err
}
