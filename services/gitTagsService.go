package services

import (
	"fmt"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/Masterminds/semver"
	"github.com/adgear/git-tags-resource/utils"
	ssh2 "golang.org/x/crypto/ssh"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GitTagsService interface
type GitTagsService interface {
	FetchTags(uri string, repositoryName string, tagFilter string) ([]string, error)
	ExtractTags(latestOnly bool, tagList []semver.Version) ([]map[string]string, error)
	CloneRef(uri string, destination string, version string) error
	CheckoutATag(destination string, version string) (map[string]string, error)
	CheckoutLWTag(destination string, version string) (map[string]string, error)
}

type gitTagsService struct {
	sshAuth *ssh.PublicKeys
}

// NewGitTagsService instance
func NewGitTagsService(keyString string, keyPassword string) (GitTagsService, error) {
	auth, err := loadPubKey(keyString, keyPassword)

	if err != nil {
		return nil, err
	}

	return gitTagsService{
		sshAuth: auth,
	}, nil
}

func loadPubKey(keyString string, keyPassword string) (*ssh.PublicKeys, error) {
	var (
		auth *ssh.PublicKeys
		err  error
	)

	if keyString != "" {
		auth, err = ssh.NewPublicKeys("git", []byte(keyString), keyPassword)
		auth.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
	}

	return auth, err
}

func (g gitTagsService) CloneRef(uri string, destination string, version string) error {
	var err error
	ref := "refs/tags/" + version
	refName := plumbing.ReferenceName(ref)

	if g.sshAuth != nil {

		_, err = git.PlainClone(destination, false, &git.CloneOptions{
			URL:               uri,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			ReferenceName:     refName,
			Auth:              g.sshAuth,
		})
	} else {

		_, err = git.PlainClone(destination, false, &git.CloneOptions{
			URL:               uri,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			ReferenceName:     refName,
		})
	}

	if err != nil {
		utils.HandleError(err)
		return err
	}

	return nil
}

func (g gitTagsService) CheckoutLWTag(destination string, version string) (map[string]string, error) {
	var info map[string]string
	var err error

	r, err := git.PlainOpen(destination)

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	w, err := r.Worktree()

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + version),
	})

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	_r, err := git.PlainOpen(destination)

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	LWTag, err := _r.Head()

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	commitObj, err := r.CommitObject(LWTag.Hash())

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	info = map[string]string{
		"tag":       version,
		"ref":       LWTag.Hash().String(),
		"shortref":  LWTag.Hash().String()[0:8],
		"tag_type":  "lw",
		"committer": commitObj.Committer.String(),
		"author":    commitObj.Author.String(),
	}

	return info, err
}

func (g gitTagsService) CheckoutATag(destination string, version string) (map[string]string, error) {
	var info map[string]string

	ref := "refs/tags/" + version
	refName := plumbing.ReferenceName(ref)

	r, err := git.PlainOpen(destination)

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	remoteRef, err := r.Reference(refName, true)

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	commitHash := remoteRef.Hash()

	aTag, err := r.TagObject(commitHash)

	if err == nil {

		commitObj, err := aTag.Commit()

		if err != nil {
			utils.HandleError(err)
			return nil, err
		}

		info = map[string]string{
			"tag":       version,
			"ref":       aTag.Hash.String(),
			"shortref":  aTag.Hash.String()[0:8],
			"tag_type":  "a",
			"committer": commitObj.Committer.String(),
			"author":    commitObj.Author.String(),
		}
	}

	if err == plumbing.ErrObjectNotFound {
		return nil, &utils.ErrATagNotFound{}
	}

	if err == plumbing.ErrReferenceNotFound {
		fmt.Println(ref)
		return nil, err
	}

	return info, nil
}

func (g gitTagsService) ExtractTags(latestOnly bool, tagList []semver.Version) ([]map[string]string, error) {
	refs := []map[string]string{}
	if latestOnly {
		refs = append(refs, map[string]string{"ref": tagList[len(tagList)-1].Original()})
	} else {
		for _, v := range tagList {
			refs = append(refs, map[string]string{"ref": v.Original()})
		}
	}

	return refs, nil
}

func (g gitTagsService) FetchTags(uri string, repositoryName string, tagFilter string) ([]string, error) {
	r, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{uri},
	})

	if err != nil {
		return nil, err
	}
	var tagRefs []*plumbing.Reference

	if g.sshAuth != nil {
		tagRefs, err = remote.List(&git.ListOptions{
			Auth: g.sshAuth,
		})

		if err != nil {
			utils.HandleError(err)
			return nil, err
		}

	} else {
		tagRefs, err = remote.List(nil)

		if err != nil {
			return nil, err
		}

	}

	if err != nil {
		utils.HandleError(err)
		return nil, err
	}

	var tagList []string

	for _, ref := range tagRefs {
		b, err := filepath.Match(tagFilter, ref.Name().Short())

		if err != nil {
			return nil, err
		}

		if b {
			tagList = append(tagList, ref.Name().Short())
		}
	}

	return tagList, nil
}
