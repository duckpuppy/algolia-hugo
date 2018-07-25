// +build mage

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/duckpuppy/algolia-hugo"
)

var ldflags = "-w -s" +
	" -X $PACKAGE/cmd.Version=$VERSION" +
	" -X $PACKAGE/cmd.Build=$BUILD_DATE" +
	" -X $PACKAGE/cmd.Commit=$COMMIT_HASH" +
	" -X $PACKAGE/cmd.Branch=$BRANCH"

var pkgPrefixLen = len(packageName)

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"
var goget = sh.RunCmd(goexe, "get", "-u")

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
		goget = sh.RunCmd(goexe, "get", "-u")
	}
}

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	branch, _ := sh.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	version := "0.0.1" // TODO: This needs to be read dynamically somehow
	return map[string]string{
		"VERSION":     version,
		"PACKAGE":     packageName,
		"BRANCH":      branch,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

func getDep() error {
	return goget("github.com/golang/dep/cmd/dep")
}

func getMetalinter() error {
	return goget("github.com/alecthomas/gometalinter")
}

func getGoreleaser() error {
	return goget("github.com/goreleaser/goreleaser")
}

func installMetalinterLinters() error {
	return sh.Run("gometalinter", "--install")
}

func getPackages() ([]string, error) {
	s, err := sh.Output(goexe, "list", "./...")
	if err != nil {
		return nil, err
	}

	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
	}
	return pkgs, nil
}

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
// nolint: deadcode
func Build() error {
	mg.Deps(Vendor)
	fmt.Println("Building...")
	return sh.RunWith(flagEnv(), goexe, "build", "-ldflags", ldflags, "-race", "-o", "build/algolia-hugo", "-v")
}

// Install Go Dep and sync vendored dependencies
// nolint: deadcode
func Vendor() error {
	mg.Deps(getDep)
	fmt.Println("Vendoring Dependencies...")
	return sh.RunV("dep", "ensure")
}

// A custom install step if you need your bin someplace other than go/bin
// nolint: deadcode
func Install() error {
	mg.Deps(Vendor)
	fmt.Println("Installing...")
	return sh.RunWith(flagEnv(), goexe, "install", "-ldflags", ldflags, "-race", "./...")
}

// Clean up after yourself
// nolint: deadcode
func Clean() {
	fmt.Println("Cleaning...")
	sh.Run(goexe, "clean", "-x")
	sh.Rm("build")
	sh.Rm("dist")
	sh.Rm("coverage.out")
	sh.Rm("coverage-all.out")
}

// Run tests
// nolint: deadcode
func Test() error {
	fmt.Println("Running Tests...")
	const (
		coverAll = "coverage-all.out"
		cover    = "coverage.out"
	)

	f, err := os.Create(coverAll)
	if err != nil {
		return err
	}

	defer f.Close()
	if _, err := f.Write([]byte("mode: count")); err != nil {
		return err
	}

	pkgs, err := getPackages()
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if err := sh.Run(goexe, "test", "-v", "-coverprofile="+cover, "-covermode=count", pkg); err != nil {
			return err
		}

		b, err := ioutil.ReadFile(cover)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		// Strip off the first line, which is the mode line
		idx := bytes.Index(b, []byte{'\n'})
		b = b[idx+1:]
		if _, err := f.Write(b); err != nil {
			return err
		}
	}

	if err := f.Close(); err != nil {
		return err
	}

	return sh.Run(goexe, "tool", "cover", "-html="+coverAll, "-o", "coverage.html")
}

// Lint using Metalinter
// nolint: deadcode
func Lint() {
	mg.Deps(getMetalinter, installMetalinterLinters)
	fmt.Println("Linting...")
	if err := sh.RunV("gometalinter", "./..."); err != nil {
		fmt.Println(err)
	}
}
