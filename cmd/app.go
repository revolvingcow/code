package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	EnvironmentVariables = map[string]string{
		"CODE_VCS": "git;hg;tf;bzr",

		"CODE_GIT_ADD":      "add",
		"CODE_GIT_CHECK":    "branch",
		"CODE_GIT_INCOMING": "log ..@{u}",
		"CODE_GIT_MERGE":    "merge",
		"CODE_GIT_MV":       "mv",
		"CODE_GIT_PULL":     "fetch",
		"CODE_GIT_PUSH":     "push",
		"CODE_GIT_RM":       "rm",
		"CODE_GIT_UPDATE":   "pull",

		"CODE_HG_ADD":      "add",
		"CODE_HG_CHECK":    "branch",
		"CODE_HG_INCOMING": "incoming",
		"CODE_HG_MERGE":    "merge",
		"CODE_HG_MV":       "mv",
		"CODE_HG_PULL":     "pull",
		"CODE_HG_PUSH":     "push",
		"CODE_HG_RM":       "rm",
		"CODE_HG_UPDATE":   "pull -u",

		"CODE_BZR_ADD":      "add",
		"CODE_BZR_CHECK":    "root",
		"CODE_BZR_INCOMING": "missing",
		"CODE_BZR_MERGE":    "merge",
		"CODE_BZR_MV":       "mv",
		"CODE_BZR_PULL":     "pull",
		"CODE_BZR_PUSH":     "push",
		"CODE_BZR_RM":       "rm",
		"CODE_BZR_UPDATE":   "update",

		"CODE_TF_ADD":      "add",
		"CODE_TF_CHECK":    "branches .",
		"CODE_TF_INCOMING": "history -r -stopafter:1 -version:W~T .",
		"CODE_TF_MERGE":    "merge",
		"CODE_TF_MV":       "rename",
		"CODE_TF_PULL":     "get -preview",
		"CODE_TF_PUSH":     "checkin",
		"CODE_TF_RM":       "delete",
		"CODE_TF_UPDATE":   "get",
	}
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
		Directory: GetWorkingDirectory(),
	}
}

// Run the application
func (a *App) Run() error {
	ConfigureEnvironment()

	if len(a.Args) < 1 {
		a.Args = append(a.Args, "help")
	}

	executed := false
	for _, s := range GetVersionControlSystems() {
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
func GetWorkingDirectory() string {
	d, _ := os.Getwd()
	return d
}

// Extract an array of version control systems available on the system
func GetVersionControlSystems() []string {
	return strings.Split(os.Getenv("CODE_VCS"), ";")
}

func ConfigureEnvironment() {
	for key, value := range EnvironmentVariables {
		env := os.Getenv(key)

		if env == "" && value != "" {
			os.Setenv(key, value)
		}
	}
}

// Determine if the directory is part of the version control system
func isVersionControlled(vcs, directory string) error {
	env := os.Getenv(fmt.Sprintf("CODE_%s_CHECK", strings.ToUpper(vcs)))
	if env == "" {
		return errors.New(fmt.Sprintf("CODE_%s_CHECK is not set", strings.ToUpper(vcs)))
	}

	err := os.Chdir(directory)
	if err != nil {
		return err
	}

	// Execute the command and swallow any output
	//var out bytes.Buffer
	actions := strings.Split(env, " ")
	cmd := exec.Command(vcs, actions...)
	//cmd.Stdout = &out

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
