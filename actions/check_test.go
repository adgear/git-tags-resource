package actions

import (
	"errors"
	"testing"

	"github.com/Masterminds/semver"

	"github.com/adgear/git-tags-resource/utils"
	"github.com/stretchr/testify/assert"
)

func TestCheckPublic_latest(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}
	output := "[{\"refs\":\"4.2.1\"}]"
	tagList := []string{"4.0.0", "4.2.1"}
	tagsVersions, _ := utils.ConvertMatchingToSemver(tagList, source.TagFilter)
	tagsResults := []map[string]string{
		map[string]string{"refs": "4.2.1"},
	}

	gtsMock.EXPECT().FetchTags(source.URI, source.RepositoryName, source.TagFilter).Return(tagList, nil).Times(1)
	gtsMock.EXPECT().ExtractTags(true, tagsVersions).Return(tagsResults, nil).Times(1)

	cr, _ := NewCheckResource(gtsMock)

	o, err := cr.Execute(source)

	assert.Equal(t, output, o, "they should be equal")

	assert.NoError(t, err)
}

func TestCheckPublic_invalidSemver(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
		LatestOnly:     "false",
	}
	v, _ := semver.NewVersion("4.2.1")
	output := []semver.Version{
		*v,
	}
	tagList := []string{"4.0.0MEH", "4.2.1"}
	tagListSemver, err := utils.ConvertMatchingToSemver(tagList, source.TagFilter)

	assert.NoError(t, err)
	assert.Equal(t, output, tagListSemver)
}

func TestCheckPublic_all(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
		LatestOnly:     "false",
	}
	output := "[{\"refs\":\"4.2.1\"},{\"refs\":\"4.0.0\"}]"
	tagList := []string{"4.0.0", "4.2.1"}
	tagsVersions, _ := utils.ConvertMatchingToSemver(tagList, source.TagFilter)
	tagsResults := []map[string]string{
		map[string]string{"refs": "4.2.1"},
		map[string]string{"refs": "4.0.0"},
	}

	gtsMock.EXPECT().FetchTags(source.URI, source.RepositoryName, source.TagFilter).Return(tagList, nil).Times(1)
	gtsMock.EXPECT().ExtractTags(false, tagsVersions).Return(tagsResults, nil).Times(1)

	cr, _ := NewCheckResource(gtsMock)

	o, err := cr.Execute(source)

	assert.Equal(t, output, o, "they should be equal")

	assert.NoError(t, err)
}

func TestCheckPublic_NoRepoName(t *testing.T) {
	setup(t)

	source := utils.Source{
		URI:       "https://github.com/concourse/concourse.git",
		TagFilter: "*.*.*",
	}

	cr, _ := NewCheckResource(gtsMock)

	_, err := cr.Execute(source)

	assert.Errorf(t, err, "repository_name can't be empty.")
}

func TestCheckPublic_GitNoKey(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "git@github.com:concourse/concourse.git",
		TagFilter:      "*.*.*",
	}

	cr, _ := NewCheckResource(gtsMock)

	_, err := cr.Execute(source)

	assert.Error(t, err, "private_key is required for git repository over SSH.")
}

func TestCheckPublic_NoTags(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}
	output := ""
	tagList := []string{}

	gtsMock.EXPECT().FetchTags(source.URI, source.RepositoryName, source.TagFilter).Return(tagList, nil).Times(1)

	cr, _ := NewCheckResource(gtsMock)

	o, err := cr.Execute(source)

	assert.Equal(t, output, o, "they should be equal")

	assert.NoError(t, err)
}

func TestCheckPublic_TagsError(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}

	tagList := []string{}

	gtsMock.EXPECT().FetchTags(source.URI, source.RepositoryName, source.TagFilter).Return(tagList, errors.New("Oops")).Times(1)

	cr, _ := NewCheckResource(gtsMock)

	_, err := cr.Execute(source)

	assert.Error(t, err)
}

func TestCheckPublic_ExtractError(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}
	tagList := []string{"4.0.0", "4.2.1"}
	tagsVersions, _ := utils.ConvertMatchingToSemver(tagList, source.TagFilter)
	tagsResults := []map[string]string{
		map[string]string{"refs": "4.2.1"},
	}

	gtsMock.EXPECT().FetchTags(source.URI, source.RepositoryName, source.TagFilter).Return(tagList, nil).Times(1)
	gtsMock.EXPECT().ExtractTags(true, tagsVersions).Return(tagsResults, errors.New("Oops")).Times(1)

	cr, _ := NewCheckResource(gtsMock)

	_, err := cr.Execute(source)

	assert.Error(t, err)
}

func TestCheckPublic_DefaultValues(t *testing.T) {
	setup(t)

	source := utils.Source{
		URI:            "git@github.com:concourse/concourse.git",
		RepositoryName: "concourse/concourse",
		PrivateKey:     "---",
		TagFilter:      "*",
	}
	tagList := []string{"4.0.0", "4.2.1"}
	tagsVersions, _ := utils.ConvertMatchingToSemver(tagList, source.TagFilter)
	tagsResults := []map[string]string{
		map[string]string{"refs": "4.2.1"},
	}

	gtsMock.EXPECT().FetchTags(source.URI, source.RepositoryName, source.TagFilter).Return(tagList, nil).Times(1)
	gtsMock.EXPECT().ExtractTags(true, tagsVersions).Return(tagsResults, nil).Times(1)

	cr, _ := NewCheckResource(gtsMock)

	_, err := cr.Execute(source)

	assert.NoError(t, err)
}
