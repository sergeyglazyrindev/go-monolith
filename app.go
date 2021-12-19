package gomonolith

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/sergeyglazyrindev/go-monolith/admin" // needed to initialize populatetemplatecontext function
	abtestblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/abtest"
	approvalblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/approval"
	authblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/auth"
	languageblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/language"
	logblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/logging"
	sessionsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/sessions"
	settingsblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/settings"
	userblueprint "github.com/sergeyglazyrindev/go-monolith/blueprint/user"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"io/fs"
	nethttp "net/http"
	"os"
	"path"
	"strconv"
)

type App struct {
	Config              *core.Config
	Database            *core.Database
	Router              *gin.Engine
	CommandRegistry     *core.CommandRegistry
	BlueprintRegistry   core.IBlueprintRegistry
	DashboardAdminPanel *core.DashboardAdminPanel
	RouterInitialized   bool
}

func (a *App) GetConfig() *core.Config {
	return a.Config
}

func (a *App) GetDatabase() *core.Database {
	return a.Database
}

func (a *App) GetRouter() *gin.Engine {
	return a.Router
}

func (a *App) GetCommandRegistry() *core.CommandRegistry {
	return a.CommandRegistry
}

func (a *App) GetBlueprintRegistry() core.IBlueprintRegistry {
	return a.BlueprintRegistry
}

func (a *App) GetDashboardAdminPanel() *core.DashboardAdminPanel {
	return a.DashboardAdminPanel
}

var appInstance *App

func NewApp(environment string, dontInitialize ...bool) *App {
	if appInstance == nil {
		a := App{}
		a.DashboardAdminPanel = core.NewDashboardAdminPanel()
		core.CurrentDashboardAdminPanel = a.DashboardAdminPanel
		a.Config = core.NewConfig("configs/" + environment + ".yml")
		a.Config.TemplatesFS = templatesRoot
		a.Config.LocalizationFS = localizationRoot
		a.CommandRegistry = &core.CommandRegistry{
			Actions: make(map[string]core.ICommand),
		}
		core.CurrentDatabaseSettings = &core.DatabaseSettings{
			Default: a.Config.D.Db.Default,
		}
		if a.Config.D.Db.Slave != nil && a.Config.D.Db.Slave.Name != "" {
			core.CurrentDatabaseSettings.Slave = a.Config.D.Db.Slave
		}
		a.BlueprintRegistry = core.NewBlueprintRegistry()
		a.Database = core.NewDatabase(a.Config)
		a.Router = gin.Default()
		//a.Router.Use(cors.New(cors.Config{
		//	AllowOrigins:     []string{"https://foo.com"},
		//	AllowMethods:     []string{"PUT", "PATCH"},
		//	AllowHeaders:     []string{"Origin"},
		//	ExposeHeaders:    []string{"Content-Length"},
		//	AllowCredentials: true,
		//	AllowOriginFunc: func(origin string) bool {
		//		return origin == "https://github.com"
		//	},
		//	MaxAge: 12 * time.Hour,
		//}))
		appInstance = &a
		a.RegisterBaseBlueprints()
		a.RegisterBaseCommands()
		if len(dontInitialize) == 0 || dontInitialize[0] == false {
			a.Initialize()
			a.InitializeRouter()
		}
		return &a
	}
	return appInstance
}

func ClearApp() {
	appInstance = nil
}

func StoreCurrentApp(app *App) {
	appInstance = app
}

func (a *App) Initialize() {
	a.BlueprintRegistry.Initialize(a)
}

func (a *App) GetAuthAdapterRegistry() *core.AuthProviderRegistry {
	return authblueprint.ConcreteBlueprint.AuthAdapterRegistry
}

func (a *App) RegisterBaseBlueprints() {
	a.BlueprintRegistry.Register(userblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(sessionsblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(settingsblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(logblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(languageblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(approvalblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(abtestblueprint.ConcreteBlueprint)
	a.BlueprintRegistry.Register(authblueprint.ConcreteBlueprint)
}

func (a *App) RegisterBlueprint(blueprint core.IBlueprint) {
	a.BlueprintRegistry.Register(blueprint)
}

func (a *App) RegisterCommand(name string, command core.ICommand) {
	a.CommandRegistry.AddAction(name, command)
}

func (a *App) RegisterBaseCommands() {
	a.RegisterCommand("migrate", &MigrateCommand{})
	a.RegisterCommand("blueprint", &BlueprintCommand{})
	a.RegisterCommand("swagger", &SwaggerCommand{})
	a.RegisterCommand("openapi", &OpenAPICommand{})
	a.RegisterCommand("superuser", &SuperadminCommand{})
	a.RegisterCommand("admin", &AdminCommand{})
	a.RegisterCommand("contenttype", &ContentTypeCommand{})
	a.RegisterCommand("generate-fake-data", &CreateFakedDataCommand{})
	a.RegisterCommand("language", &LanguageCommand{})
	a.RegisterCommand("dbshell", &DbShellCommand{})
	a.RegisterCommand("shell", &ShellCommand{})
	a.RegisterCommand("microservice-example", &MicroserviceExampleCommand{})
}

func (a *App) ExecuteCommand() {
	var action string
	var isCorrectActionPassed bool = false
	var help string
	if len(os.Args) > 1 {
		action = os.Args[1]
		isCorrectActionPassed = a.CommandRegistry.IsRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := a.CommandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return
	}
	if len(os.Args) > 2 {
		subaction := os.Args[2]
		isCorrectActionPassed = a.CommandRegistry.IsRegisteredCommand(action)
		err := a.CommandRegistry.RunAction(action, subaction, os.Args[3:])
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := a.CommandRegistry.RunAction(action, "", make([]string, 0))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (a *App) TriggerCommandExecution(action string, subaction string, params []string) {
	a.CommandRegistry.RunAction(action, subaction, params)
}

func (a *App) StartAdmin() {
	// useradmin.RegisterAdminPart()
}

func (a *App) StartAPI() {
	a.Initialize()
	// _ = a.Router.Run(":" + strconv.Itoa(a.Config.D.API.ListenPort))
}

//go:embed templates
var templatesRoot embed.FS

//go:embed localization
var localizationRoot embed.FS

//go:embed static/*
var staticRoot embed.FS

// myFS implements fs.FS
type goMonolithStaticFS struct {
	content embed.FS
}

func (c goMonolithStaticFS) Open(name string) (fs.File, error) {
	return c.content.Open(path.Join("static", name))
}

func (a *App) InitializeRouter() {
	if a.RouterInitialized {
		return
	}
	// http.FS can be used to create a http Filesystem
	staticFiles := goMonolithStaticFS{staticRoot}
	fs1 := nethttp.FS(staticFiles)
	// Serve static files
	a.Router.StaticFS("/static-inbuilt/", fs1)
	a.BlueprintRegistry.InitializeRouting(a, a.Router)
	a.DashboardAdminPanel.RegisterHTTPHandlers(a.Router)
	a.RouterInitialized = true
}

func (a *App) BaseAPIUrl() string {
	return ":" + strconv.Itoa(a.Config.D.API.ListenPort)
}
