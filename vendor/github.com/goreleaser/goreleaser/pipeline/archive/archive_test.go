package archive

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

func TestRunPipe(t *testing.T) {
	var assert = assert.New(t)
	folder, err := ioutil.TempDir("", "archivetest")
	assert.NoError(err)
	current, err := os.Getwd()
	assert.NoError(err)
	assert.NoError(os.Chdir(folder))
	defer func() {
		assert.NoError(os.Chdir(current))
	}()
	var dist = filepath.Join(folder, "dist")
	assert.NoError(os.Mkdir(dist, 0755))
	assert.NoError(os.Mkdir(filepath.Join(dist, "mybin_darwin_amd64"), 0755))
	assert.NoError(os.Mkdir(filepath.Join(dist, "mybin_windows_amd64"), 0755))
	_, err = os.Create(filepath.Join(dist, "mybin_darwin_amd64", "mybin"))
	assert.NoError(err)
	_, err = os.Create(filepath.Join(dist, "mybin_windows_amd64", "mybin.exe"))
	assert.NoError(err)
	_, err = os.Create(filepath.Join(folder, "README.md"))
	assert.NoError(err)
	var ctx = &context.Context{
		Folders: map[string]string{
			"darwinamd64":  "mybin_darwin_amd64",
			"windowsamd64": "mybin_windows_amd64",
		},
		Config: config.Project{
			Dist: dist,
			Archive: config.Archive{
				Files: []string{
					"README.*",
				},
				FormatOverrides: []config.FormatOverride{
					{
						Goos:   "windows",
						Format: "zip",
					},
				},
			},
		},
	}
	for _, format := range []string{"tar.gz", "zip"} {
		t.Run("Archive format "+format, func(t *testing.T) {
			ctx.Config.Archive.Format = format
			assert.NoError(Pipe{}.Run(ctx))
		})
	}
}

func TestRunPipeBinary(t *testing.T) {
	var assert = assert.New(t)
	folder, err := ioutil.TempDir("", "archivetest")
	assert.NoError(err)
	current, err := os.Getwd()
	assert.NoError(err)
	assert.NoError(os.Chdir(folder))
	defer func() {
		assert.NoError(os.Chdir(current))
	}()
	var dist = filepath.Join(folder, "dist")
	assert.NoError(os.Mkdir(dist, 0755))
	assert.NoError(os.Mkdir(filepath.Join(dist, "mybin_darwin"), 0755))
	assert.NoError(os.Mkdir(filepath.Join(dist, "mybin_win"), 0755))
	_, err = os.Create(filepath.Join(dist, "mybin_darwin", "mybin"))
	assert.NoError(err)
	_, err = os.Create(filepath.Join(dist, "mybin_win", "mybin.exe"))
	assert.NoError(err)
	_, err = os.Create(filepath.Join(folder, "README.md"))
	assert.NoError(err)
	var ctx = &context.Context{
		Folders: map[string]string{
			"darwinamd64":  "mybin_darwin",
			"windowsamd64": "mybin_win",
		},
		Config: config.Project{
			Dist: dist,
			Builds: []config.Build{
				{Binary: "mybin"},
			},
			Archive: config.Archive{
				Format: "binary",
			},
		},
	}
	assert.NoError(Pipe{}.Run(ctx))
	assert.Contains(ctx.Artifacts, "mybin_darwin/mybin")
	assert.Contains(ctx.Artifacts, "mybin_win/mybin.exe")
	assert.Len(ctx.Artifacts, 2)
}

func TestRunPipeDistRemoved(t *testing.T) {
	var assert = assert.New(t)
	var ctx = &context.Context{
		Folders: map[string]string{
			"darwinamd64": "whatever",
		},
		Config: config.Project{
			Dist: "/path/nope",
			Archive: config.Archive{
				Format: "zip",
			},
		},
	}
	assert.Error(Pipe{}.Run(ctx))
}

func TestRunPipeInvalidGlob(t *testing.T) {
	var assert = assert.New(t)
	var ctx = &context.Context{
		Folders: map[string]string{
			"windowsamd64": "whatever",
		},
		Config: config.Project{
			Dist: "/tmp",
			Archive: config.Archive{
				Files: []string{
					"[x-]",
				},
			},
		},
	}
	assert.Error(Pipe{}.Run(ctx))
}

func TestRunPipeGlobFailsToAdd(t *testing.T) {
	var assert = assert.New(t)
	folder, err := ioutil.TempDir("", "archivetest")
	assert.NoError(err)
	current, err := os.Getwd()
	assert.NoError(err)
	assert.NoError(os.Chdir(folder))
	defer func() {
		assert.NoError(os.Chdir(current))
	}()
	assert.NoError(os.MkdirAll(filepath.Join(folder, "folder", "another"), 0755))

	var ctx = &context.Context{
		Folders: map[string]string{
			"windows386": "mybin",
		},
		Config: config.Project{
			Dist: folder,
			Archive: config.Archive{
				Files: []string{
					"folder",
				},
			},
		},
	}
	assert.Error(Pipe{}.Run(ctx))
}
