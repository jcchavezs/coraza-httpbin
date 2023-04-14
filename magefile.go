//go:build mage
// +build mage

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

var golangCILintVer = "v1.48.0" // https://github.com/golangci/golangci-lint/releases
var gosImportsVer = "v0.1.5"    // https://github.com/rinchsan/gosimports/releases/tag/v0.1.5

var errRunGoModTidy = errors.New("go.mod/sum not formatted, commit changes")

// Lint verifies code quality.
func Lint() error {
	if err := sh.RunV("go", "run", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", golangCILintVer), "run"); err != nil {
		return err
	}

	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}

	if sh.Run("git", "diff", "--exit-code", "go.mod", "go.sum") != nil {
		return errRunGoModTidy
	}

	return nil
}

func Test() error {
	return sh.RunV("go", "test", "-v", "./...")
}

func build(goos string, goarch string) error {
	if err := os.MkdirAll("build", 0755); err != nil {
		return err
	}

	suffix := ""
	env := map[string]string{}
	if goos != "" {
		suffix += "-" + goos
		env["GOOS"] = goos
	}

	if goarch != "" {
		suffix += "-" + goarch
		env["GOARCH"] = goarch
	}

	return sh.RunWithV(env, "go", "build", "-o", "build/coraza-httpbin"+suffix, "cmd/coraza-httpbin/main.go")
}

// Build builds the project
func Build() error {
	return build("", "")
}

func BuildForDockerImage() error {
	if err := build("linux", "amd64"); err != nil {
		return err
	}

	if err := build("linux", "arm64"); err != nil {
		return err
	}

	return nil
}

const dockerImage = "ghcr.io/jcchavezs/coraza-httpbin"

func PackDockerImage() error {
	if err := BuildForDockerImage(); err != nil {
		return err
	}

	if err := sh.RunV("docker", "buildx", "build", "-t", dockerImage, "--platform=linux/arm64,linux/amd64,darwin/amd64", "."); err != nil {
		return err
	}

	return nil
}

func PackLocalDockerImage() error {
	if err := BuildForDockerImage(); err != nil {
		return err
	}

	if err := sh.RunWithV(map[string]string{"DOCKER_BUILDKIT": "1"}, "docker", "buildx", "build", "-t", dockerImage, "--load", "."); err != nil {
		return err
	}

	return nil
}
