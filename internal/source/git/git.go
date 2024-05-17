package git

import (
	"gcmerge/internal/flags"
	"os"
	"reflect"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type GitT struct {
	ConfigFilepath string
	StoredConfig   []byte
	IsUpdated      bool

	SshKeyFilepath string
	RepoSshUrl     string
	RepoPath       string
	RepoBranch     string
}

func (s *GitT) Init(f flags.DaemonFlagsT) {
	s.SshKeyFilepath = f.GitSshKeyFilepath
	s.RepoSshUrl = f.GitSshUrl
	s.RepoBranch = f.GitBranch
	s.RepoPath = f.TmpDir + "/repo"
	s.ConfigFilepath = f.SourcePath
}

func (s *GitT) GetConfig() (config []byte, err error) {
	s.IsUpdated = false

	if _, err = os.Stat(s.RepoPath); !os.IsNotExist(err) {
		if err = os.RemoveAll(s.RepoPath); err != nil {
			return config, err
		}
	}

	if _, err = os.Stat(s.SshKeyFilepath); err != nil {
		return config, err
	}

	publicSshKey, err := ssh.NewPublicKeysFromFile("git", s.SshKeyFilepath, "")
	if err != nil {
		return config, err
	}

	_, err = git.PlainClone(s.RepoPath, false, &git.CloneOptions{
		URL:           s.RepoSshUrl,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(s.RepoBranch),
		SingleBranch:  true,
		Auth:          publicSshKey,
	})
	if err != nil {
		return config, err
	}

	if config, err = os.ReadFile(s.ConfigFilepath); err != nil {
		return config, err
	}

	if !reflect.DeepEqual(s.StoredConfig, config) {
		s.IsUpdated = true
		s.StoredConfig = config
	}

	return config, err
}

func (s *GitT) NeedUpdate() bool {
	return s.IsUpdated
}
