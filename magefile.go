// +build mage

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
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

// func getDep() error {
// 	return goget("github.com/golang/dep/cmd/dep")
// }

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

// Build the binary
func Build() error { // nolint: deadcode
	mg.Deps(Vendor)
	fmt.Println("Building...")
	return sh.RunWith(flagEnv(), goexe, "build", "-ldflags", ldflags, "-race", "-o", "build/algolia-hugo", "-v")
}

// Install Go Dep and sync vendored dependencies
func Vendor() error { // nolint: deadcode
	// mg.Deps(getDep)
	fmt.Println("Vendoring Dependencies...")
	return sh.RunV("dep", "ensure")
}

// Install the source and binary
func Install() error { // nolint: deadcode
	mg.Deps(Vendor)
	fmt.Println("Installing...")
	return sh.RunWith(flagEnv(), goexe, "install", "-ldflags", ldflags, "-race", "./...")
}

// Clean up after yourself
func Clean() { // nolint: deadcode
	fmt.Println("Cleaning...")
	sh.Run(goexe, "clean", "-x")
	sh.Rm("build")
	sh.Rm("dist")
	sh.Rm("coverage.out")
	sh.Rm("coverage-all.out")
}

// Run tests
func Test() error { // nolint: deadcode
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
func Lint() { // nolint: deadcode
	// We have to do these sequentially rather than in one call
	mg.Deps(getMetalinter)
	mg.Deps(installMetalinterLinters)

	fmt.Println("Linting...")
	if err := sh.RunV("gometalinter", "./..."); err != nil {
		fmt.Println(err)
	}
}

// Run all tests and linters
func Check() { // nolint: deadcode
	if strings.Contains(runtime.Version(), "1.8") {
		fmt.Println("Skip Check on %s\n", runtime.Version())
		return
	}

	mg.Deps(Test, Lint)
}
