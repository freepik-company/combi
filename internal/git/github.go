package git

import (
	"github.com/google/go-github/v61/github"
)

func NewGitHubClient() (client *github.Client) {
	client = github.NewClient(nil)
	return client
}
