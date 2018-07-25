package fpm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/goreleaser/goreleaser/config"
	"github.com/goreleaser/goreleaser/context"
	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	assert.NotEmpty(t, Pipe{}.Description())
}

func TestRunPipeNoFormats(t *testing.T) {
	var assert = assert.New(t)
	var ctx = &context.Context{
		Config: config.Project{},
	}
	assert.NoError(Pipe{}.Run(ctx))
}

func TestRunPipeFormatBinary(t *testing.T) {
	var assert = assert.New(t)
	var ctx = &context.Context{
		Config: config.Project{
			FPM: config.FPM{
				Formats: []string{"deb"},
			},
			Archive: config.Archive{
				Format: "binary",
			},
		},
	}
	assert.NoError(Pipe{}.Run(ctx))
}

func TestRunPipe(t *testing.T) {
	var assert = assert.New(t)
	folder, err := ioutil.TempDir("", "archivetest")
	assert.NoError(err)
	var dist = filepath.Join(folder, "dist")
	assert.NoError(os.Mkdir(dist, 0755))
	assert.NoError(os.Mkdir(filepath.Join(dist, "mybin"), 0755))
	_, err = os.Create(filepath.Join(dist, "mybin", "mybin"))
	assert.NoError(err)
	var ctx = &context.Context{
		Folders: map[string]string{
			"linuxamd64":  "mybin",
			"linux386":    "mybin",
			"darwinamd64": "anotherbin",
		},
		Config: config.Project{
			ProjectName: "mybin",
			Dist:        dist,
			FPM: config.FPM{
				Formats:      []string{"deb"},
				Dependencies: []string{"make"},
				Conflicts:    []string{"git"},
				Description:  "Some description",
				License:      "MIT",
				Maintainer:   "me@me",
				Vendor:       "asdf",
				Homepage:     "https://goreleaser.github.io",
			},
		},
	}
	assert.NoError(Pipe{}.Run(ctx))
}

func TestNoFPMInPath(t *testing.T) {
	var assert = assert.New(t)
	var path = os.Getenv("PATH")
	defer func() {
		assert.NoError(os.Setenv("PATH", path))
	}()
	assert.NoError(os.Setenv("PATH", ""))
	var ctx = &context.Context{
		Config: config.Project{
			FPM: config.FPM{
				Formats: []string{"deb"},
			},
		},
	}
	assert.EqualError(Pipe{}.Run(ctx), ErrNoFPM.Error())
}

func TestCreateFileDoesntExist(t *testing.T) {
	var assert = assert.New(t)
	folder, err := ioutil.TempDir("", "archivetest")
	assert.NoError(err)
	var dist = filepath.Join(folder, "dist")
	assert.NoError(os.Mkdir(dist, 0755))
	assert.NoError(os.Mkdir(filepath.Join(dist, "mybin"), 0755))
	var ctx = &context.Context{
		Folders: map[string]string{
			"linuxamd64": "mybin",
		},
		Config: config.Project{
			Dist: dist,
			FPM: config.FPM{
				Formats: []string{"deb"},
			},
		},
	}
	assert.Error(Pipe{}.Run(ctx))
}

func TestCreatePathDoesntExist(t *testing.T) {
	var assert = assert.New(t)
	var ctx = &context.Context{}
	assert.Contains(
		create(ctx, "tar.gz", "/nope/noooo", "windowsarm64").Error(),
		"no such file",
	)
}
