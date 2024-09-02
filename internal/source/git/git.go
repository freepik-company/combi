package git

import (
	"fmt"
	"os"
	"reflect"

	"combi/api/v1alpha2"
	"combi/internal/globals"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type GitSourceT struct {
	ConfigFilepath string
	StoredConfig   []byte

	SshKeyFilepath string
	RepoSshUrl     string
	RepoPath       string
	RepoBranch     string
}

func (s *GitSourceT) Init(source v1alpha2.SourceT) {

	s.SshKeyFilepath = source.Git.SshKeyFilepath
	s.RepoSshUrl = source.Git.SshUrl
	s.RepoBranch = source.Git.Branch
	s.RepoPath = globals.TmpDir + "/repo"
	s.ConfigFilepath = fmt.Sprintf("%s/%s", s.RepoPath, source.Git.Filepath)
}

func (s *GitSourceT) GetConfig() (config []byte, updated bool, err error) {
	if _, err = os.Stat(s.RepoPath); !os.IsNotExist(err) {
		if err = os.RemoveAll(s.RepoPath); err != nil {
			return config, updated, err
		}
	}

	if _, err = os.Stat(s.SshKeyFilepath); err != nil {
		return config, updated, err
	}

	publicSshKey, err := ssh.NewPublicKeysFromFile("git", s.SshKeyFilepath, "")
	if err != nil {
		return config, updated, err
	}

	_, err = git.PlainClone(s.RepoPath, false, &git.CloneOptions{
		URL:           s.RepoSshUrl,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(s.RepoBranch),
		SingleBranch:  true,
		Auth:          publicSshKey,
	})
	if err != nil {
		return config, updated, err
	}

	if config, err = os.ReadFile(s.ConfigFilepath); err != nil {
		return config, updated, err
	}
	config = []byte(os.ExpandEnv(string(config)))

	if !reflect.DeepEqual(s.StoredConfig, config) {
		updated = true
	}

	return config, updated, err
}
