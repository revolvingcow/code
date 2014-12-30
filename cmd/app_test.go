package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestExtension(t *testing.T) {
	var buffer bytes.Buffer

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

func TestRawCommand(t *testing.T) {
	var buffer bytes.Buffer
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

func TestMissingCommand(t *testing.T) {
	var buffer bytes.Buffer
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

func TestMissingRepository(t *testing.T) {
	os.Chdir(os.TempDir())

	var buffer bytes.Buffer
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
