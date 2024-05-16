package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Git struct {
	SshKeyFilepath     string
	RepoSshUrl         string
	RepoPath           string
	RepoBranch         string
	RepoConfigFilepath string
	RepoStoredHash     string
	RepoIsUpdated      bool
}

func (g *Git) GetConfig() (config []byte, err error) {
	g.RepoIsUpdated = false

	if _, err = os.Stat(g.RepoPath); !os.IsNotExist(err) {
		if err = os.RemoveAll(g.RepoPath); err != nil {
			return config, err
		}
	}

	if _, err = os.Stat(g.SshKeyFilepath); err != nil {
		return config, err
	}

	publicSshKey, err := ssh.NewPublicKeysFromFile("git", g.SshKeyFilepath, "")
	if err != nil {
		return config, err
	}

	repo, err := git.PlainClone(g.RepoPath, false, &git.CloneOptions{
		URL:           g.RepoSshUrl,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(g.RepoBranch),
		SingleBranch:  true,
		Auth:          publicSshKey,
	})
	if err != nil {
		return config, err
	}

	headRef, err := repo.Head()
	if err != nil {
		return config, err
	}

	currentHash := headRef.Hash().String()
	if g.RepoStoredHash != currentHash {
		g.RepoStoredHash = currentHash
		g.RepoIsUpdated = true

		if config, err = os.ReadFile(g.RepoPath + "/" + g.RepoConfigFilepath); err != nil {
			return config, err
		}
	}

	return config, err
}

func (g *Git) NeedUpdate() bool {
	return g.RepoIsUpdated
}
