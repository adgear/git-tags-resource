package actions

//go:generate mockgen -destination=../mocks/mock_in.go -package=mocks github.com/adgear/git-tags-resource InResource

import (
	"errors"

	"github.com/adgear/git-tags-resource/services"
	"github.com/adgear/git-tags-resource/utils"
)

// InResource interface
type InResource interface {
	Execute(source utils.Source, destination string, version string) (string, error)
}

type inResource struct {
	GTS services.GitTagsService
}

// NewInResource returns a new instance
func NewInResource(gts services.GitTagsService) (InResource, error) {
	return inResource{
		GTS: gts,
	}, nil
}

// Execute the in resource
func (ir inResource) Execute(source utils.Source, destination string, version string) (string, error) {
	if destination == "" {
		return "", errors.New("empty destination")
	}

	var (
		output string
		info   map[string]string
	)

	err := ir.GTS.CloneRef(source.URI, destination, version)
	if err != nil {
		return "", err
	}

	info, err = ir.GTS.CheckoutATag(destination, version)
	if err != nil {
		if err.Error() == "not a tag" {
			info, err = ir.GTS.CheckoutLWTag(destination, version)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	err = utils.WriteInfo(destination, info)

	if err != nil {
		return "", err
	}

	metadata, err := utils.BuildMetadata(info)

	if err != nil {
		return "", err
	}

	output = "{version: {ref: \"" + version + "\"}, metadata: " + metadata + "}"

	return output, nil
}
