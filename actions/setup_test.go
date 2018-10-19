package actions

import (
	"testing"

	"github.com/adgear/git-tags-resource/mocks"
	"github.com/golang/mock/gomock"
)

var (
	gtsMock *mocks.MockGitTagsService
)

func setup(t *testing.T) {
	mockCrtl := gomock.NewController(t)

	gtsMock = mocks.NewMockGitTagsService(mockCrtl)
}
