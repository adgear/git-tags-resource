package utils

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/Masterminds/semver"
)

// RemoveContents DANGEROUS... I lost all my code :'( freaking typos...
func RemoveContents(destination string) {
}

// BuildMetadata in right format for concourse
func BuildMetadata(info map[string]string) (string, error) {
	metadata := []map[string]string{}
	for k, v := range info {
		metadata = append(metadata, map[string]string{"name": k, "value": v})
	}

	m, err := json.Marshal(metadata)

	if err != nil {
		return "", err
	}

	return string(m), nil
}

// WriteInfo writes some valuable stuff to disk.
func WriteInfo(destination string, info map[string]string) error {

	for k, v := range info {
		err := ioutil.WriteFile(filepath.Join(destination, ".git", k), []byte(v), 0644)

		if err != nil {
			return err
		}
	}

	return nil
}

// ConvertMatchingToSemver is sorting and applying tag filter
func ConvertMatchingToSemver(tagList []string, tagFilter string) ([]semver.Version, error) {
	tagListSemver := []semver.Version{}

	for _, v := range tagList {
		match, err := regexp.MatchString("^"+tagFilter+"$", v)
		if err != nil {
			return nil, err
		}

		if match {
			ver, err := semver.NewVersion(v)
			if err != nil {
				return nil, err
			}
			tagListSemver = append(tagListSemver, *ver)
		}
	}

	return tagListSemver, nil
}

// ByVersion sort semver
type ByVersion []semver.Version

func (a ByVersion) Len() int           { return len(a) }
func (a ByVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVersion) Less(i, j int) bool { return a[i].LessThan(&a[j]) }
