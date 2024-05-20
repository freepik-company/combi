package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
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
	repoFilePath := "config/combi.yaml"

	for {
		storer := memory.NewStorage()
		workTree := memfs.New()
		repo, err := git.Clone(storer, workTree,
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

			commit, err := repo.CommitObject(headRef.Hash())
			if err != nil {
				log.Fatalf("unable to get commit in local repository: %s", err.Error())
			}

			tree, _ := commit.Tree()

			objFile, _ := tree.File(repoFilePath)

			fileStr, _ := objFile.Contents()

			fmt.Print(fileStr)
		}

		fmt.Printf("waiting to next pull...\n")
		time.Sleep(5 * time.Second)
	}
}
