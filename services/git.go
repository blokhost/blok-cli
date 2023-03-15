package services

import (
	"fmt"
	"github.com/babilu-online/common/context"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"os"
)

type GitService struct {
	context.DefaultService

	auth *AuthService
}

const GIT_SVC = "git_svc"

func (svc GitService) Id() string {
	return GIT_SVC
}

func (svc *GitService) Start() error {
	svc.auth = svc.Service(AUTH_SVC).(*AuthService)

	return nil
}

func (svc *GitService) InGitDir(path string) (bool, error) {
	path = fmt.Sprintf("%s/.git", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (svc *GitService) Download(src string, dst string) error {
	isGit, err := svc.InGitDir(dst)
	if err != nil {
		return err
	}

	if isGit {
		return svc.Pull()
	}

	_, err = git.PlainClone(dst, false, &git.CloneOptions{
		URL: src,
		Auth: &http.BasicAuth{
			Username: "blok", //can be anything except an empty string
			Password: svc.auth.GithubToken(),
		},
		Progress: os.Stdout,
	})
	return err
}

func (svc *GitService) Pull() error {
	r, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "blok", // yes, this can be anything except an empty string
			Password: svc.auth.GithubToken(),
		}})

	ref, err := r.Head()

	if err != nil {
		return err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	log.Println("Pulled Latest Hash:", commit.Hash)
	return nil
}
