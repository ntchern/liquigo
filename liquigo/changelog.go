package liquigo

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ChangelogFiles represents a list of files to apply
type changelogFiles struct {
	// Path of the changeset.
	Path string

	// Body of the changeset
	Files []string `yaml:"databaseChangeLog"`
}

func parseChangelog(r io.Reader) (changelogFiles, error) {
	result := changelogFiles{}
	yamlFile, err := ioutil.ReadAll(r)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
