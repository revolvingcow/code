package cmd

import (
	"bytes"
	"os"
	"testing"
)

// Test using an extension command
func TestExtension(t *testing.T) {
	var buffer bytes.Buffer

	os.Setenv("CODE_VCS", "git")
	os.Setenv("CODE_GIT_CHECK", "branch")
	os.Setenv("CODE_GIT_INCOMING", "log ..@{u}")

	app := App{
		Args:      []string{"incoming"},
		Stdin:     &buffer,
		Stderr:    &buffer,
		Stdout:    &buffer,
		Directory: getWorkingDirectory(),
	}

	err := app.Run()
	if err != nil {
		t.FailNow()
	}
}

// Test a raw command not currently defined but can pass through
func TestRawCommand(t *testing.T) {
	var buffer bytes.Buffer

	os.Setenv("CODE_VCS", "git")
	os.Setenv("CODE_GIT_CHECK", "branch")

	app := App{
		Args:      []string{"status"},
		Stdin:     &buffer,
		Stderr:    &buffer,
		Stdout:    &buffer,
		Directory: getWorkingDirectory(),
	}

	err := app.Run()
	if err != nil {
		t.FailNow()
	}
}

// Test behavior when a VCS "check" is not defined
func TestMissingCheck(t *testing.T) {
	var buffer bytes.Buffer

	os.Setenv("CODE_VCS", "git")
	os.Unsetenv("CODE_GIT_CHECK")

	app := App{
		Args:      []string{"status"},
		Stdin:     &buffer,
		Stderr:    &buffer,
		Stdout:    &buffer,
		Directory: getWorkingDirectory(),
	}

	err := app.Run()
	if err == nil {
		t.FailNow()
	}
}

// Test behavior when a subcommand is issued which does not exist
func TestMissingCommand(t *testing.T) {
	var buffer bytes.Buffer

	os.Setenv("CODE_VCS", "git")
	os.Setenv("CODE_GIT_CHECK", "branch")

	app := App{
		Args:      []string{"nonexistant", "command"},
		Stdin:     &buffer,
		Stderr:    &buffer,
		Stdout:    &buffer,
		Directory: getWorkingDirectory(),
	}

	err := app.Run()
	if err != nil {
		t.FailNow()
	}
}

// Test output if there is currently no VCS repository found
func TestMissingRepository(t *testing.T) {
	var buffer bytes.Buffer

	os.Chdir(os.TempDir())
	os.Setenv("CODE_VCS", "git")
	os.Setenv("CODE_GIT_CHECK", "branch")

	app := App{
		Args:      []string{"status"},
		Stdin:     &buffer,
		Stderr:    &buffer,
		Stdout:    &buffer,
		Directory: getWorkingDirectory(),
	}

	err := app.Run()
	if err == nil {
		t.FailNow()
	}
}

// Test handling of no arguments passed
func TestNoArguments(t *testing.T) {
	var buffer bytes.Buffer

	os.Setenv("CODE_VCS", "git")
	os.Setenv("CODE_GIT_CHECK", "branch")

	app := App{
		Args:      []string{},
		Stdin:     &buffer,
		Stderr:    &buffer,
		Stdout:    &buffer,
		Directory: getWorkingDirectory(),
	}

	err := app.Run()
	if err == nil {
		t.FailNow()
	}
}

// Test the creation of a new application
func TestNewApp(t *testing.T) {
	app := NewApp()

	if app.Stdin == nil {
		t.Fail()
	}

	if app.Stderr == nil {
		t.Fail()
	}

	if app.Stdout == nil {
		t.Fail()
	}

	if app.Directory == "" {
		t.Fail()
	}
}
