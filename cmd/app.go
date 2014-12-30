package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// The application
type App struct {
	Args      []string
	Stdin     io.Reader
	Stderr    io.Writer
	Stdout    io.Writer
	Directory string
}

// Create a new application instance
func NewApp() *App {
	return &App{
		Args:      os.Args[1:],
		Stdin:     os.Stdin,
		Stderr:    os.Stderr,
		Stdout:    os.Stdout,
		Directory: getWorkingDirectory(),
	}
}

// Run the application
func (a *App) Run() error {
	if len(a.Args) < 1 {
		a.Args = append(a.Args, "help")
	}

	executed := false
	for _, s := range getVersionControlSystems() {
		// Determine if the directory is version controlled (skip if it is not)
		err := isVersionControlled(s, a.Directory)
		if err != nil {
			continue
		}

		// Execute the subcommand
		err = a.executeSubcommand(s, a.Args[0], a.Args[1:]...)
		executed = true
	}

	// If nothing was executed inform the user there is no repository found
	if !executed {
		return errors.New("No repository found")
	}

	return nil
}

// Get the working directory
func getWorkingDirectory() string {
	d, _ := os.Getwd()
	return d
}

// Extract an array of version control systems available on the system
func getVersionControlSystems() []string {
	return strings.Split(os.Getenv("CODE_VCS"), ";")
}

// Determine if the directory is part of the version control system
func isVersionControlled(vcs, directory string) error {
	env := os.Getenv(fmt.Sprintf("CODE_%s_CHECK", strings.ToUpper(vcs)))
	if env == "" {
		return errors.New(fmt.Sprintf("CODE_%s_CHECK is not set", strings.ToUpper(vcs)))
	}

	// Execute the command and swallow any output
	var out bytes.Buffer
	actions := strings.Split(env, " ")
	cmd := exec.Command(vcs, actions...)
	cmd.Stdout = &out

	return cmd.Run()
}

// Execute a subcommand of the given version control system while passing along all arguments
func (a *App) executeSubcommand(vcs, subcommand string, args ...string) error {
	command := []string{}
	env := os.Getenv(fmt.Sprintf("CODE_%s_%s", strings.ToUpper(vcs), strings.ToUpper(subcommand)))

	// If there is an override to the subcommand we will utilize it
	// otherwise, we simply pass everything to the VCS in raw format
	if env == "" {
		command = append(command, subcommand)
		command = append(command, args...)
	} else {
		actions := strings.Split(env, " ")
		command = append(command, actions...)
		command = append(command, args...)
	}

	// Execute the subcommand and output everything to standard system resources
	cmd := exec.Command(vcs, command...)
	cmd.Stdout = a.Stdout
	cmd.Stderr = a.Stderr
	cmd.Stdin = a.Stdin

	return cmd.Run()
}
