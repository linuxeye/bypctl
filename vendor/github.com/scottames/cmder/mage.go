//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/scottames/cmder"
)

var (
	g0 = cmder.New("go").RunFn()
)

// Check run checks
func Check() {
	mg.Deps(Lint, Test)
}

// Lint run golangci-lint
func Lint() error {
	return cmder.New("golangci-lint", "run").Run()
}

// Test run go tests
func Test() error {
	return g0("test", "./...")
}

// Tidy run go mod tidy
func Tidy() error {
	return g0("mod", "tidy")
}
