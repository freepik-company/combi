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

	// Prepare ssh key
	_, err := os.Stat(env[gitRepoPrivateKeyKey])
	if err != nil {
		log.Fatalf("read file %s failed: %s\n", env[gitRepoPrivateKeyKey], err.Error())
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", env[gitRepoPrivateKeyKey], "")
	if err != nil {
		log.Fatalf("generate public key: %s\n", err.Error())
	}

	// Prepare cloned repo to pull
	repoPath := "/tmp/repo"
	branchReferenceName := plumbing.NewBranchReferenceName(env[gitRepoBranchKey])
	repo, err := git.PlainClone(repoPath, false,
		&git.CloneOptions{
			URL: env[gitRepoSshUrlKey],
			// Depth:         1,
			ReferenceName: branchReferenceName,
			SingleBranch:  true,
			Auth:          publicKeys,
		},
	)
	if err != nil {
		log.Fatalf("unable to clone the repository: %s", err.Error())
	}

	workTree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("unable to get repository work tree: %s", err.Error())
	}

	// Prepare HEAD commit hash
	headRef, err := repo.Head()
	if err != nil {
		log.Fatalf("unable to set HEAD in the repository: %s", err.Error())
	}

	currentHash := ""
	storedHash := headRef.Hash().String()
	fmt.Printf("the stored hash of '%s' branch is: %s\n", env[gitRepoBranchKey], storedHash)

	for {
		err = workTree.Pull(&git.PullOptions{
			SingleBranch: true,
			// Depth:         5,
			ReferenceName: branchReferenceName,
			Auth:          publicKeys,
			Force:         true,
			Progress:      os.Stdout,
		})
		if err != nil {
			if err.Error() != git.NoErrAlreadyUpToDate.Error() &&
				err.Error() != git.ErrNonFastForwardUpdate.Error() {
				log.Fatalf("unable to pull repository: %s", err.Error())
			}

			if err.Error() != git.ErrNonFastForwardUpdate.Error() {
				ref, err := repo.Reference(plumbing.NewRemoteReferenceName("origin", env[gitRepoBranchKey]), true)
				if err != nil {
					log.Fatalf("unable to get remote reference repository: %s", err.Error())
				}

				err = workTree.Reset(&git.ResetOptions{
					Mode:   git.HardReset,
					Commit: ref.Hash(),
				})
				if err != nil {
					log.Fatalf("unable to reset repository: %s", err.Error())
				}
			}
		}

		headRef, err := repo.Head()
		if err != nil {
			log.Fatalf("unable to set HEAD in the repository: %s", err.Error())
		}

		currentHash = headRef.Hash().String()
		if storedHash != currentHash {
			fmt.Printf("the HEAD hash change from '%s' to %s\n", storedHash, currentHash)
			storedHash = currentHash
		}

		fmt.Printf("waiting to next pull...\n")
		time.Sleep(5 * time.Second)
	}
}
