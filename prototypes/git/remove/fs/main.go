package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

const (
	gitRepoPrivateKeyKey = "GIT_REPO_PRIVATE_KEY"
	gitRepoSshUrlKey     = "GIT_REPO_SSH_URL"
	gitRepoBranchKey     = "GIT_REPO_BRANCH"
	gitRepoFilepathKey   = "GIT_REPO_FILEPATH"
)

var (
	env = map[string]string{
		gitRepoPrivateKeyKey: "",
		gitRepoSshUrlKey:     "",
		gitRepoBranchKey:     "",
		gitRepoFilepathKey:   "",
	}
)

func main() {
	for key := range env {
		env[key] = os.Getenv(key)
		if env[key] == "" {
			log.Fatalf("env var %s not provided", key)
		}
	}

	_, err := os.Stat(env[gitRepoPrivateKeyKey])
	if err != nil {
		log.Fatalf("read file %s failed: %s\n", env[gitRepoPrivateKeyKey], err.Error())
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", env[gitRepoPrivateKeyKey], "")
	if err != nil {
		log.Fatalf("generate public key: %s\n", err.Error())
	}

	storedHash := ""
	currentHash := ""
	repoPath := "/tmp/repo"
	repoFilePath := "config/gcmerge.yaml"

	for {
		repo, err := git.PlainClone(repoPath, false,
			&git.CloneOptions{
				URL:           env[gitRepoSshUrlKey],
				Depth:         1,
				ReferenceName: plumbing.NewBranchReferenceName(env[gitRepoBranchKey]),
				SingleBranch:  true,
				Auth:          publicKeys,
			},
		)
		if err != nil {
			log.Fatalf("unable to clone the repository: %s", err.Error())
		}

		headRef, err := repo.Head()
		if err != nil {
			log.Fatalf("unable to set HEAD in the repository: %s", err.Error())
		}

		currentHash = headRef.Hash().String()
		if storedHash != currentHash {
			fmt.Printf("the HEAD hash change from '%s' to %s\n", storedHash, currentHash)

			storedHash = currentHash

			file, err := os.ReadFile(repoPath + "/" + repoFilePath)
			if err != nil {
				log.Fatalf("unable to get config file in local repository: %s", err.Error())
			}
			fmt.Print(string(file))
		}

		if err = os.RemoveAll(repoPath); err != nil {
			log.Fatalf("unable to remove local repository: %s", err.Error())
		}

		fmt.Printf("waiting to next pull...\n")
		time.Sleep(5 * time.Second)
	}
}
