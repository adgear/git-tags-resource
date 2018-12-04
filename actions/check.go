package actions

//go:generate mockgen -destination=../mocks/mock_check.go -package=mocks github.com/adgear/git-tags-resource CheckResource

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/adgear/git-tags-resource/services"
	"github.com/adgear/git-tags-resource/utils"
)

// CheckResource interface
type CheckResource interface {
	Execute(source utils.Source) (string, error)
}

type checkResource struct {
	GTS services.GitTagsService
}

// NewCheckResource returns a new instance
func NewCheckResource(gts services.GitTagsService) (CheckResource, error) {
	return checkResource{
		GTS: gts,
	}, nil
}

// Execute the check resource
func (cr checkResource) Execute(source utils.Source) (string, error) {
	valid := true
	msg := ""

	if source.RepositoryName == "" {
		msg = "repository_name can't be empty."
		valid = false
	}

	if source.URI == "" {
		source.URI = "git@github.com:" + source.RepositoryName + ".git"
	}

	if strings.HasPrefix(source.URI, "git") && source.PrivateKey == "" {
		msg = "private_key is required for git repository over SSH."
		valid = false
	}

	if strings.HasPrefix(source.URI, "http") {
		source.PrivateKey = ""
		source.PrivateKeyPassword = ""
	}

	if source.TagFilter == "" {
		source.TagFilter = "*"
	}

	if source.LatestOnly == "" {
		source.LatestOnly = "true"
	}

	var refs []map[string]string

	if valid {
		var err error

		tagList, err := cr.GTS.FetchTags(source.URI, source.RepositoryName, source.TagFilter)

		if err != nil {
			return "", err
		}

		if len(tagList) == 0 {
			return "", nil
		}

		sort.Strings(tagList)

		s, err := utils.ConvertMatchingToSemver(tagList, source.TagFilter)

		if err != nil {
			return "", err
		}

		sort.Sort(utils.ByVersion(s))

		latestOnly, err := strconv.ParseBool(source.LatestOnly)

		if err != nil {
			return "", err
		}

		refs, err = cr.GTS.ExtractTags(latestOnly, s)

		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New(msg)
	}

	output, err := json.Marshal(refs)

	if err != nil {
		return "", nil
	}

	return string(output), nil
}
