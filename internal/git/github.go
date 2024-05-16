package git

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func CloneSshRepo(localRepoPath string, publicSshKey *ssh.PublicKeys) (err error) {
	return err
}

func getSshKey(privateKeyFilepath string) (publicSshKey *ssh.PublicKeys, err error) {
	_, err = os.Stat(privateKeyFilepath)
	if err != nil {
		return publicSshKey, err
	}

	publicSshKey, err = ssh.NewPublicKeysFromFile("git", privateKeyFilepath, "")
	if err != nil {
		return publicSshKey, err
	}

	return publicSshKey, err
}
