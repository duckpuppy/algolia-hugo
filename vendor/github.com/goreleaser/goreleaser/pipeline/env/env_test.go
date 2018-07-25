package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/goreleaser/goreleaser/config"
	"github.com/goreleaser/goreleaser/context"
	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	assert.NotEmpty(t, Pipe{}.Description())
}

func TestValidEnv(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(os.Setenv("GITHUB_TOKEN", "asdf"))
	var ctx = &context.Context{
		Config:   config.Project{},
		Validate: true,
		Publish:  true,
	}
	assert.NoError(Pipe{}.Run(ctx))
}

func TestInvalidEnv(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(os.Unsetenv("GITHUB_TOKEN"))
	var ctx = &context.Context{
		Config:   config.Project{},
		Validate: true,
		Publish:  true,
	}
	assert.Error(Pipe{}.Run(ctx))
}

type flags struct {
	Validate, Publish, Snapshot bool
}

func TestInvalidEnvChecksSkipped(t *testing.T) {
	for _, flag := range []flags{
		{
			Validate: false,
			Publish:  true,
		}, {
			Validate: true,
			Publish:  false,
		}, {
			Validate: true,
		},
	} {
		t.Run(fmt.Sprintf("%v", flag), func(t *testing.T) {
			var assert = assert.New(t)
			assert.NoError(os.Unsetenv("GITHUB_TOKEN"))
			var ctx = &context.Context{
				Config:   config.Project{},
				Validate: flag.Validate,
				Publish:  flag.Publish,
				Snapshot: flag.Snapshot,
			}
			assert.NoError(Pipe{}.Run(ctx))
		})
	}
}
