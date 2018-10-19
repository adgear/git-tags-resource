package actions

import (
	"errors"
	"testing"

	"github.com/adgear/git-tags-resource/utils"
	"github.com/stretchr/testify/assert"
)

func TestInPublic_ATag(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}

	input := utils.Input{
		Source:  source,
		Version: map[string]string{"refs": "4.2.1"},
	}
	destination := "/tmp/test.xxx"

	output := "{version: {ref: \"4.2.1\"}, metadata: []}"

	gtsMock.EXPECT().CloneRef(input.Source.URI, destination, input.Version["refs"]).Return(nil).Times(1)
	gtsMock.EXPECT().CheckoutATag(destination, input.Version["refs"]).Return(map[string]string{}, nil).Times(1)

	ir, _ := NewInResource(gtsMock)

	o, err := ir.Execute(input.Source, destination, input.Version["refs"])

	assert.Equal(t, output, o)

	assert.NoError(t, err)
}

func TestInPublic_LWTag(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}

	input := utils.Input{
		Source:  source,
		Version: map[string]string{"refs": "4.2.1"},
	}
	destination := "/tmp/test.xxx"

	output := "{version: {ref: \"4.2.1\"}, metadata: []}"

	gtsMock.EXPECT().CloneRef(input.Source.URI, destination, input.Version["refs"]).Return(nil).Times(1)
	gtsMock.EXPECT().CheckoutATag(destination, input.Version["refs"]).Return(map[string]string{}, errors.New("not a tag")).Times(1)
	gtsMock.EXPECT().CheckoutLWTag(destination, input.Version["refs"]).Return(map[string]string{}, nil).Times(1)

	ir, _ := NewInResource(gtsMock)

	o, err := ir.Execute(input.Source, destination, input.Version["refs"])

	assert.Equal(t, output, o)

	assert.NoError(t, err)
}

func TestInPublic_NoDestination(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}

	input := utils.Input{
		Source:  source,
		Version: map[string]string{"refs": "4.2.1"},
	}
	destination := ""

	ir, _ := NewInResource(gtsMock)

	_, err := ir.Execute(input.Source, destination, input.Version["refs"])

	assert.EqualError(t, err, "empty destination")
}

func TestInPublic_CloneRefError(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}

	input := utils.Input{
		Source:  source,
		Version: map[string]string{"refs": "4.2.1"},
	}
	destination := "/tmp/test.xxx"

	gtsMock.EXPECT().CloneRef(input.Source.URI, destination, input.Version["refs"]).Return(errors.New("")).Times(1)

	ir, _ := NewInResource(gtsMock)

	_, err := ir.Execute(input.Source, destination, input.Version["refs"])

	assert.Error(t, err)
}

func TestInPublic_ATagError(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}

	input := utils.Input{
		Source:  source,
		Version: map[string]string{"refs": "4.2.1"},
	}
	destination := "/tmp/test.xxx"

	gtsMock.EXPECT().CloneRef(input.Source.URI, destination, input.Version["refs"]).Return(nil).Times(1)
	gtsMock.EXPECT().CheckoutATag(destination, input.Version["refs"]).Return(map[string]string{}, errors.New("some other error")).Times(1)

	ir, _ := NewInResource(gtsMock)

	_, err := ir.Execute(input.Source, destination, input.Version["refs"])

	assert.Error(t, err)
}

func TestInPublic_LWTagError(t *testing.T) {
	setup(t)

	source := utils.Source{
		RepositoryName: "concourse/concourse",
		URI:            "https://github.com/concourse/concourse.git",
		TagFilter:      "*.*.*",
	}

	input := utils.Input{
		Source:  source,
		Version: map[string]string{"refs": "4.2.1"},
	}
	destination := "/tmp/test.xxx"

	gtsMock.EXPECT().CloneRef(input.Source.URI, destination, input.Version["refs"]).Return(nil).Times(1)
	gtsMock.EXPECT().CheckoutATag(destination, input.Version["refs"]).Return(map[string]string{}, errors.New("not a tag")).Times(1)
	gtsMock.EXPECT().CheckoutLWTag(destination, input.Version["refs"]).Return(map[string]string{}, errors.New("some other error")).Times(1)

	ir, _ := NewInResource(gtsMock)

	_, err := ir.Execute(input.Source, destination, input.Version["refs"])

	assert.Error(t, err)
}
