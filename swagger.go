package gomonolith

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"io"
	"log"
	"os"
	"os/exec"
)

type SwaggerCommand struct {
}

func (c SwaggerCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &core.CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}

	commandRegistry.AddAction("serve", &ServeSwaggerServer{})
	if len(os.Args) > 2 {
		action = os.Args[2]
		isCorrectActionPassed = commandRegistry.IsRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return nil
	}
	return commandRegistry.RunAction(subaction, "", args)
}

func (c SwaggerCommand) GetHelpText() string {
	return "Manage your swagger integration"
}

type ServeSwaggerServerOptions struct {
}

type ServeSwaggerServer struct {
}

func (command ServeSwaggerServer) Proceed(subaction string, args []string) error {
	appInstance.Config.APISpec = core.NewSwaggerSpec(appInstance.Config.D.Swagger.PathToSpec)
	commandToExecute := exec.Command(
		"swagger", "serve", "--flavor=swagger", "--no-open",
		fmt.Sprintf("--port=%d", appInstance.Config.D.Swagger.ListenPort), appInstance.Config.D.Swagger.PathToSpec,
	)
	stderr, err := commandToExecute.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err := commandToExecute.Start(); err != nil {
		log.Fatal(err)
		return err
	}

	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := commandToExecute.Wait(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (command ServeSwaggerServer) GetHelpText() string {
	return "Serve your swagger api spec"
}
