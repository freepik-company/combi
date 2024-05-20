package source

import (
	"combi/internal/flags"
	"combi/internal/source/git"
	"combi/internal/source/local"
)

type SourceT interface {
	Init(f flags.DaemonFlagsT)
	GetConfig() (config []byte, err error)
	NeedUpdate() bool
}

func GetSources() (sources map[string]SourceT) {
	sources = map[string]SourceT{
		"local": &local.LocalT{},
		"git":   &git.GitT{},
	}
	return sources
}
