package gomonolith

import (
	_ "github.com/traefik/yaegi" // this is our go shell package
	"os"
	"os/exec"
)

type ShellCommand struct {
}

func (c ShellCommand) Proceed(subaction string, args []string) error {
	commandToExecute := exec.Command(
		"yaegi",
	)
	// Sets standard output to cmd.stdout writer
	commandToExecute.Stdout = os.Stdout
	// Sets standard input to cmd.stdin reader
	commandToExecute.Stdin = os.Stdin
	var err error
	if err = commandToExecute.Start(); err != nil {
		return err
	}
	if err = commandToExecute.Wait(); err != nil {
		return err
	}
	return nil
}

func (c ShellCommand) GetHelpText() string {
	return "Start gore go shell"
}
