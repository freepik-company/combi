package source

import (
	"combi/api/v1alpha2"
	"combi/internal/source/file"
	"combi/internal/source/git"
	"combi/internal/source/kube"
	"combi/internal/source/raw"
)

type SourceT interface {
	Init(source v1alpha2.SourceT)
	GetConfig() (config []byte, updated bool, err error)
}

func GetSources() (sources map[string]SourceT) {
	sources = map[string]SourceT{
		"raw":        &raw.RawSourceT{},
		"file":       &file.FileSourceT{},
		"git":        &git.GitSourceT{},
		"kubernetes": &kube.KubeSourceT{},
	}
	return sources
}
