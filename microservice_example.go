package gomonolith

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	interfaces3 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/interfaces"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"os"
)

type MicroserviceExample struct {
	core.Microservice
}

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
type ValidationError struct {
	Code    string
	Message string
	Params  []string
}

// A GeneralError is an error that is used when something strange happened
// swagger:response generalError
type GeneralError struct {
	Error string
}

type MicroserviceExampleCommand struct {
}

func (c MicroserviceExampleCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &core.CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}

	commandRegistry.AddAction("start", &MicroserviceExampleStartCommand{})
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

func (c MicroserviceExampleCommand) GetHelpText() string {
	return "Start your example microservice"
}

type MicroserviceExampleStartOptions struct {
	StartSwagger bool `long:"swagger" description:"Start swagger for your microservice"`
}

type MicroserviceExampleStartCommand struct {
}

func (command MicroserviceExampleStartCommand) Proceed(subaction string, args []string) error {
	var opts = &MicroserviceExampleStartOptions{}
	parser := flags.NewParser(opts, flags.Default)
	var err error
	_, err = parser.ParseArgs(args)
	if err != nil {
		return err
	}
	appInstance.GetAuthAdapterRegistry().RegisterNewAdapter(&interfaces3.TokenAuthProvider{})
	microservice := &MicroserviceExample{Microservice: core.Microservice{
		Port: 8089, AuthBackend: "token", Name: "Example microservice",
		Prefix: "ExampleMicroservice", SwaggerPort: 8090, ServiceSwaggerDefinition: &core.ServiceSwaggerDefinition{
			BasePath: "/v3", Host: "localhost", Schemes: []string{"http", "https"}, Info: &core.ServiceSwaggerDefinitionInfo{
				Version: "0.0.1",
				Contact: &core.ServiceSwaggerDefinitionInfoContact{
					Name:  "Sergey Glazyrin",
					URL:   "https://github.com/sergeyglazyrindev/go-monolith",
					Email: "sergey.glazyrin.dev@gmail.com",
				},
				License: &core.ServiceSwaggerDefinitionInfoLicense{
					Name: "MIT",
					URL:  "http://opensource.org/licenses/MIT",
				},
			},
			Consumes:       []string{"application/json", "application/xml"},
			Produces:       []string{"application/json", "application/xml"},
			SwaggerVersion: "2.0",
		},
	}}
	if opts.StartSwagger {
		return microservice.StartSwagger(appInstance)
	}
	r := microservice.RegisterEndpoints(appInstance)
	microservice.Start(r)
	return nil
}

func (command MicroserviceExampleStartCommand) GetHelpText() string {
	return "Start your microservice"
}
