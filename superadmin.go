package gomonolith

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jessevdk/go-flags"
	"github.com/miquella/ask"
	utils2 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/utils"
	userblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/user"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"os"
)

type SuperadminCommand struct {
}

func (c SuperadminCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &core.CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}
	commandRegistry.AddAction("create", &CreateSuperadmin{})
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

func (c SuperadminCommand) GetHelpText() string {
	return "Manage your superusers"
}

type SuperadminCommandOptions struct {
	Username  string `short:"n" required:"true" description:"Username" valid:"username-go-monolith,username-unique"`
	Email     string `short:"e" required:"true" description:"Email'" valid:"email,email-unique"`
	FirstName string `short:"f" required:"false" description:"First name'"`
	LastName  string `short:"l" required:"false" description:"Last name'"`
}

type CreateSuperadmin struct {
}

func (command CreateSuperadmin) Proceed(subaction string, args []string) error {
	var opts = &SuperadminCommandOptions{}
	parser := flags.NewParser(opts, flags.Default)
	var err error
	_, err = parser.ParseArgs(args)
	if len(args) == 0 {
		var help string = `
Please provide flags -n and -e which are username and email of the user respectively 
`
		fmt.Printf(help)
		return nil
	}
	if err != nil {
		return err
	}
	_, err = govalidator.ValidateStruct(opts)
	if err != nil {
		return err
	}
	database := core.NewDatabaseInstance()
	defer database.Close()
	db := database.Db
	if opts.FirstName == "" {
		opts.FirstName = "System"
	}
	if opts.LastName == "" {
		opts.LastName = "Admin"
	}
	err = ask.Print("Warning! I am about to ask you for a password!\n")
	if err != nil {
		return err
	}
	var password string
	for true {
		password, err = ask.HiddenAsk("Password: ")
		if err != nil {
			return err
		}
		confirmpassword, err := ask.HiddenAsk("Confirm Password: ")
		if err != nil {
			return err
		}
		passwordValidationStruct := &userblueprint.PasswordValidationStruct{
			Password:          password,
			ConfirmedPassword: confirmpassword,
		}
		_, err = govalidator.ValidateStruct(passwordValidationStruct)
		if err != nil {
			core.Trail(core.ERROR, errors.New("please try to to repeat password again"))
			continue
		}
		break
	}
	salt := core.GenerateRandomString(appInstance.Config.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, err := utils2.HashPass(password, salt)
	if err != nil {
		return err
	}
	admin := core.GenerateUserModel()
	admin.SetFirstName(opts.FirstName)
	admin.SetLastName(opts.LastName)
	admin.SetUsername(opts.Username)
	admin.SetEmail(opts.Email)
	admin.SetPassword(hashedPassword)
	admin.SetActive(true)
	admin.SetIsSuperUser(true)
	admin.SetSalt(salt)
	admin.SetIsPasswordUsable(true)
	db.Create(admin)
	core.Trail(core.INFO, "Superuser created successfully")
	return nil
}

func (command CreateSuperadmin) GetHelpText() string {
	return "Create superadmin"
}
