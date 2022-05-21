//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/sh"
	"github.com/princjef/mageutil/bintool"
)

const binDir = "tools"

var (
	linter = bintool.Must(bintool.New(
		"golangci-lint{{.BinExt}}",
		"1.45.2",
		"https://github.com/golangci/golangci-lint/releases/download/v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
		bintool.WithFolder(binDir),
	))

	godoc = bintool.Must(bintool.NewGo("godoc", "latest", bintool.WithFolder("tools")))
)

func Lint() error {
	if err := linter.Ensure(); err != nil {
		return err
	}

	return linter.Command(`run`).Run()
}

func LintFix() error {
	if err := linter.Ensure(); err != nil {
		return err
	}

	return linter.Command(`run --fix`).Run()
}

func Test() error {
	return sh.RunV("go", "test", "-race", "-v", "-count", "1", "./...")
}

func Vendor() error {
	return sh.RunV("go", "mod", "vendor")
}

func Todo() error {
	return sh.RunV("rg", "-i", "todo")
}

func Godoc() error {
	if err := godoc.Ensure(); err != nil {
		return err
	}

	return godoc.Command("-http=:6060").Run()
}
