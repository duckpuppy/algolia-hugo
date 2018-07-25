// Package archive implements the pipe interface with the intent of
// archiving and compressing the binaries, readme, and other artifacts. It
// also provides an Archive interface which represents an archiving format.
package archive

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/goreleaser/archive"
	"github.com/goreleaser/goreleaser/context"
	"github.com/goreleaser/goreleaser/internal/archiveformat"
	"github.com/mattn/go-zglob"
	"golang.org/x/sync/errgroup"
)

// Pipe for archive
type Pipe struct{}

// Description of the pipe
func (Pipe) Description() string {
	return "Creating archives"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	var g errgroup.Group
	for platform, folder := range ctx.Folders {
		folder := folder
		platform := platform
		g.Go(func() error {
			if ctx.Config.Archive.Format == "binary" {
				return skip(ctx, platform, folder)
			}
			return create(ctx, platform, folder)
		})
	}
	return g.Wait()
}

func create(ctx *context.Context, platform, name string) error {
	var folder = filepath.Join(ctx.Config.Dist, name)
	var format = archiveformat.For(ctx, platform)
	file, err := os.Create(folder + "." + format)
	if err != nil {
		return err
	}
	log.WithField("archive", file.Name()).Info("creating")
	defer func() { _ = file.Close() }()
	var archive = archive.New(file)
	defer func() { _ = archive.Close() }()

	files, err := findFiles(ctx)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err = archive.Add(f, f); err != nil {
			return err
		}
	}
	var basepath = filepath.Join(ctx.Config.Dist, name)
	binaries, err := ioutil.ReadDir(basepath)
	if err != nil {
		return err
	}
	for _, binary := range binaries {
		var path = filepath.Join(basepath, binary.Name())
		if err := archive.Add(binary.Name(), path); err != nil {
			return err
		}
	}
	ctx.AddArtifact(file.Name())
	return nil
}

func skip(ctx *context.Context, platform, name string) error {
	var path = filepath.Join(ctx.Config.Dist, name)
	binaries, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	log.WithField("platform", platform).Debugf("found %v binaries", len(binaries))
	for _, binary := range binaries {
		log.WithField("binary", binary.Name()).Info("skip archiving")
		ctx.AddArtifact(filepath.Join(path+"/", binary.Name()))
	}
	return nil
}

func findFiles(ctx *context.Context) (result []string, err error) {
	for _, glob := range ctx.Config.Archive.Files {
		files, err := zglob.Glob(glob)
		if err != nil {
			return result, err
		}
		result = append(result, files...)
	}
	return
}
