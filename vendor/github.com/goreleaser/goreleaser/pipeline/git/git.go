// Package git implements the Pipe interface getting and validating the
// current git repository state
package git

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/context"
)

// ErrInvalidVersionFormat is return when the version isnt in a valid format
type ErrInvalidVersionFormat struct {
	version string
}

func (e ErrInvalidVersionFormat) Error() string {
	return fmt.Sprintf("%v is not in a valid version format", e.version)
}

// ErrDirty happens when the repo has uncommitted/unstashed changes
type ErrDirty struct {
	status string
}

func (e ErrDirty) Error() string {
	return fmt.Sprintf("git is currently in a dirty state:\n%v", e.status)
}

// ErrWrongRef happens when the HEAD reference is different from the tag being built
type ErrWrongRef struct {
	commit, tag string
}

func (e ErrWrongRef) Error() string {
	return fmt.Sprintf("git tag %v was not made against commit %v", e.tag, e.commit)
}

// ErrNoTag happens if the underlying git repository doesn't contain any tags
// but no snapshot-release was requested.
var ErrNoTag = fmt.Errorf("git doesn't contain any tags. Either add a tag or use --snapshot")

// Pipe for brew deployment
type Pipe struct{}

// Description of the pipe
func (Pipe) Description() string {
	return "Getting and validating git state"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) (err error) {
	tag, commit, err := getInfo()
	if err != nil {
		return
	}
	if tag == "" && !ctx.Snapshot {
		return ErrNoTag
	}
	ctx.Git = context.GitInfo{
		CurrentTag: tag,
		Commit:     commit,
	}
	if err = setLog(ctx, tag, commit); err != nil {
		return
	}
	if err = setVersion(ctx, tag, commit); err != nil {
		return
	}
	if !ctx.Validate {
		log.Warn("skipped validations because --skip-validate is set")
		return nil
	}
	return validate(ctx, commit, tag)
}

func setVersion(ctx *context.Context, tag, commit string) (err error) {
	if ctx.Snapshot {
		snapshotName, err := getSnapshotName(ctx, tag, commit)
		if err != nil {
			return fmt.Errorf("failed to generate snapshot name: %s", err.Error())
		}
		ctx.Version = snapshotName
		return nil
	}
	// removes usual `v` prefix
	ctx.Version = strings.TrimPrefix(tag, "v")
	return
}

func setLog(ctx *context.Context, tag, commit string) (err error) {
	if ctx.ReleaseNotes != "" {
		return
	}
	var log string
	if tag == "" {
		log, err = getChangelog(commit)
	} else {
		log, err = getChangelog(tag)
	}
	if err != nil {
		return err
	}
	ctx.ReleaseNotes = fmt.Sprintf("## Changelog\n\n%v", log)
	return nil
}

type snapshotNameData struct {
	Commit    string
	Tag       string
	Timestamp int64
}

func getSnapshotName(ctx *context.Context, tag, commit string) (string, error) {
	tmpl, err := template.New("snapshot").Parse(ctx.Config.Snapshot.NameTemplate)
	var out bytes.Buffer
	if err != nil {
		return "", err
	}
	var data = snapshotNameData{
		Commit:    commit,
		Tag:       tag,
		Timestamp: time.Now().Unix(),
	}
	err = tmpl.Execute(&out, data)
	return out.String(), err
}

func validate(ctx *context.Context, commit, tag string) error {
	out, err := git("status", "--porcelain")
	if strings.TrimSpace(out) != "" || err != nil {
		return ErrDirty{out}
	}
	if ctx.Snapshot {
		return nil
	}
	if !regexp.MustCompile("^[0-9.]+").MatchString(ctx.Version) {
		return ErrInvalidVersionFormat{ctx.Version}
	}
	_, err = cleanGit("describe", "--exact-match", "--tags", "--match", tag)
	if err != nil {
		return ErrWrongRef{commit, tag}
	}
	return nil
}

func getChangelog(tag string) (string, error) {
	prev, err := previous(tag)
	if err != nil {
		return "", err
	}
	if !prev.Tag {
		return gitLog(prev.SHA, tag)
	}
	return gitLog(fmt.Sprintf("%v..%v", prev.SHA, tag))
}

func gitLog(refs ...string) (string, error) {
	var args = []string{"log", "--pretty=oneline", "--abbrev-commit"}
	args = append(args, refs...)
	return git(args...)
}

func getInfo() (tag, commit string, err error) {
	tag, err = cleanGit("describe", "--tags", "--abbrev=0")
	if err != nil {
		log.WithError(err).Info("failed to retrieve current tag")
	}
	commit, err = cleanGit("show", "--format='%H'", "HEAD")
	return
}

func previous(tag string) (result ref, err error) {
	result.Tag = true
	result.SHA, err = cleanGit("describe", "--tags", "--abbrev=0", tag+"^")
	if err != nil {
		result.Tag = false
		result.SHA, err = cleanGit("rev-list", "--max-parents=0", "HEAD")
	}
	return
}

type ref struct {
	Tag bool
	SHA string
}
