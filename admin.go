package gomonolith

import (
	"fmt"
	"github.com/miquella/ask"
	"github.com/sergeyglazyrindev/go-monolith/colors"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"log"
	"os"
	"strings"
)

type AdminCommand struct {
}

func (c AdminCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &core.CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}

	commandRegistry.AddAction("serve", &ServeAdminServer{})
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
	commandRegistry.RunAction(subaction, "", args)
	return nil
}

func (c AdminCommand) GetHelpText() string {
	return "Admin functionality for Go-Monolith project"
}

type AdminStartServerOptions struct {
}

type ServeAdminServer struct {
}

const welcomeMessage = "" +
	`         ___       __          _` + "\n" +
	colors.FGBlueB + `  __  __` + colors.FGNormal + `/   | ____/ /___ ___  (_)___` + "\n" +
	colors.FGBlueB + ` / / / /` + colors.FGNormal + ` /| |/ __  / __ '__ \/ / __ \` + "\n" +
	colors.FGBlueB + `/ /_/ /` + colors.FGNormal + ` ___ / /_/ / / / / / / / / / /` + "\n" +
	colors.FGBlueB + `\__,_/` + colors.FGNormal + `_/  |_\__,_/_/ /_/ /_/_/_/ /_/` + "\n"

func (command ServeAdminServer) Proceed(subaction string, args []string) error {
	appInstance.InitializeRouter()
	migrateCommand := MigrateCommand{}
	err := migrateCommand.Proceed("determine-conflicts", []string{})
	if err != nil {
		core.Trail(core.CRITICAL, "Found problems with migrations")
		err = ask.Print("Warning! Found problems with migrations.\n")
		if err != nil {
			return err
		}
		var answer string
		for true {
			answer, err = ask.HiddenAsk("Do you want to start server ?")
			if err != nil {
				return err
			}
			if !core.Contains([]string{"yes", "no"}, strings.ToLower(answer)) {
				continue
			}
			break
		}
		if answer == "no" {
			core.Trail(core.WARNING, "You decided to solve first migration problems, so see you next time!")
			return nil
		}
	}
	core.Trail(core.OK, "Server Started: http://%s:%d", core.CurrentConfig.D.Admin.BindIP, core.CurrentConfig.D.Admin.ListenPort)
	fmt.Println(welcomeMessage)
	log.Println(appInstance.Router.Run(fmt.Sprintf("%s:%d", core.CurrentConfig.D.Admin.BindIP, core.CurrentConfig.D.Admin.ListenPort)))
	return nil
}

func (command ServeAdminServer) GetHelpText() string {
	return "Serve your admin panel"
}
