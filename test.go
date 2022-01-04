package gomonolith

import (
	"errors"
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/auth"
	interfaces3 "github.com/sergeyglazyrindev/go-monolith/blueprint/auth/interfaces"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	App      *App
	Database *core.ProjectDatabase
}

func (suite *TestSuite) SetupTest() {
	app := NewFullAppForTests()
	suite.App = app
}

func (suite *TestSuite) StoreDatabase(database *core.ProjectDatabase) {
	suite.Database = database
}

func (suite *TestSuite) TearDownSuite() {
	ClearTestApp()
}

func ClearTestApp() {
	appForTests = nil
	appInstance = nil
}

func failOnPanic(t *testing.T) {
	r := recover()
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

func newSuiteInformation() *suite.SuiteInformation {
	testStats := make(map[string]*suite.TestInformation)

	return &suite.SuiteInformation{
		TestStats: testStats,
	}
}

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }

func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Test", name); !ok {
		return false, nil
	}
	return regexp.MatchString("", name)
}

func startStats(s *suite.SuiteInformation, testName string) {
	s.TestStats[testName] = &suite.TestInformation{
		TestName: testName,
		Start:    time.Now(),
	}
}

func endStats(s *suite.SuiteInformation, testName string, passed bool) {
	s.TestStats[testName].End = time.Now()
	s.TestStats[testName].Passed = passed
}

func RunTests(t *testing.T, currentsuite suite.TestingSuite) {
	core.ProjectEmailSenderFactory.MakeEmailSender = func() core.IEmailSender {
		return core.TestEmailSender
	}
	defer failOnPanic(t)

	currentsuite.SetT(t)

	var suiteSetupDone bool

	var stats *suite.SuiteInformation
	if _, ok := currentsuite.(suite.WithStats); ok {
		stats = newSuiteInformation()
	}

	tests := []testing.InternalTest{}
	methodFinder := reflect.TypeOf(currentsuite)
	suiteName := methodFinder.Elem().Name()

	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)

		ok, err := methodFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}

		if !ok {
			continue
		}

		if !suiteSetupDone {
			if stats != nil {
				stats.Start = time.Now()
			}

			if setupAllSuite, ok := currentsuite.(suite.SetupAllSuite); ok {
				setupAllSuite.SetupSuite()
			}

			suiteSetupDone = true
		}

		test := testing.InternalTest{
			Name: method.Name,
			F: func(t *testing.T) {
				parentT := currentsuite.T()
				currentsuite.SetT(t)
				defer failOnPanic(t)
				defer func() {
					if stats != nil {
						passed := !t.Failed()
						endStats(stats, method.Name, passed)
					}

					if afterTestSuite, ok := currentsuite.(suite.AfterTest); ok {
						afterTestSuite.AfterTest(suiteName, method.Name)
					}

					if tearDownTestSuite, ok := currentsuite.(suite.TearDownTestSuite); ok {
						tearDownTestSuite.TearDownTest()
					}

					currentsuite.SetT(parentT)
				}()

				core.TestEmailSender.ClearTestEmails()
				config := core.NewConfig("configs/" + os.Getenv("TEST_ENVIRONMENT") + ".yml")
				core.CurrentConfig = config
				core.CurrentConfig.InTests = true
				core.CurrentConfig.TemplatesFS = templatesRoot
				core.CurrentConfig.LocalizationFS = localizationRoot
				core.CurrentDatabaseSettings = &core.DatabaseSettings{
					Default: config.D.Db.Default,
					Slave:   config.D.Db.Slave,
				}
				if config.D.Db.Default.Type == "sqlite" {
					a := NewApp(os.Getenv("TEST_ENVIRONMENT"), true)
					a.Config.InTests = true
					core.CurrentConfig.InTests = true
					database := core.NewDatabaseInstance()
					database.Adapter.SetTimeZone(database.Db, "UTC")
					core.TestInstanceDatabase = database
					reflect.ValueOf(currentsuite).MethodByName("StoreDatabase").Call([]reflect.Value{reflect.ValueOf(database)})
					if setupTestSuite, ok := currentsuite.(suite.SetupTestSuite); ok {
						setupTestSuite.SetupTest()
					}
					upCommand := MigrateCommand{}
					upCommand.Proceed("up", make([]string, 0))
					if beforeTestSuite, ok := currentsuite.(suite.BeforeTest); ok {
						beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
					}

					if stats != nil {
						startStats(stats, method.Name)
					}
					method.Func.Call([]reflect.Value{reflect.ValueOf(currentsuite)})
					appInstance.BlueprintRegistry.ResetMigrationTree()
					downCommand := MigrateCommand{}
					downCommand.Proceed("down", make([]string, 0))
					//err := os.Remove(suite.app.Config.D.Db.Default.Name)
					//if err != nil {
					//	assert.Equal(suite.T(), true, false, fmt.Errorf("Couldnt remove db with name %s", suite.app.Config.D.Db.Default.Name))
					//}
					core.TestInstanceDatabase = nil
					database.Close()
				} else {
					a := NewApp(os.Getenv("TEST_ENVIRONMENT"), true)
					core.CurrentConfig.InTests = true
					a.Config.InTests = true
					if !CreatedDatabaseForTests && core.CurrentDatabaseSettings.Default.Type != "sqlite" {
						adapter := core.NewDbAdapter(nil, core.CurrentDatabaseSettings.Default.Type)
						adapter.InitializeDatabaseForTests(core.CurrentDatabaseSettings.Default)
						CreatedDatabaseForTests = true
					}
					database := core.NewDatabaseInstance()
					core.TestInstanceDatabase = database
					database.Adapter.SetIsolationLevelForTests(database.Db)
					database.Adapter.SetTimeZone(database.Db, "UTC")
					if setupTestSuite, ok := currentsuite.(suite.SetupTestSuite); ok {
						setupTestSuite.SetupTest()
					}
					upCommand := MigrateCommand{}
					upCommand.Proceed("up", make([]string, 0))
					database.Db.Transaction(func(tx *gorm.DB) error {
						newDatabase := &core.ProjectDatabase{Db: tx, Adapter: database.Adapter}
						core.TestInstanceDatabase = newDatabase
						reflect.ValueOf(currentsuite).MethodByName("StoreDatabase").Call([]reflect.Value{reflect.ValueOf(newDatabase)})
						if setupTestSuite, ok := currentsuite.(suite.SetupTestSuite); ok {
							setupTestSuite.SetupTest()
						}
						if beforeTestSuite, ok := currentsuite.(suite.BeforeTest); ok {
							beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
						}

						if stats != nil {
							startStats(stats, method.Name)
						}
						method.Func.Call([]reflect.Value{reflect.ValueOf(currentsuite)})
						// return nil will commit the whole transaction
						return errors.New("dont commit")
					})
					core.TestInstanceDatabase.Adapter.ClearTestDatabase()
					core.TestInstanceDatabase = nil
					database.Close()
				}
			},
		}
		tests = append(tests, test)
	}
	if suiteSetupDone {
		defer func() {
			if tearDownAllSuite, ok := currentsuite.(suite.TearDownAllSuite); ok {
				tearDownAllSuite.TearDownSuite()
			}

			if suiteWithStats, measureStats := currentsuite.(suite.WithStats); measureStats {
				stats.End = time.Now()
				suiteWithStats.HandleStats(suiteName, stats)
			}
		}()
	}

	runTests(t, tests)
}

type runner interface {
	Run(name string, f func(t *testing.T)) bool
}

func runTests(t testing.TB, tests []testing.InternalTest) {
	if len(tests) == 0 {
		t.Log("warning: no tests to run")
		return
	}

	r, ok := t.(runner)
	if !ok { // backwards compatibility with Go 1.6 and below
		if !testing.RunTests(allTestsFilter, tests) {
			t.Fail()
		}
		return
	}

	for _, test := range tests {
		r.Run(test.Name, test.F)
	}
}

//func NewTestApp() *App {
//	a := App{}
//	a.DashboardAdminPanel = core.NewDashboardAdminPanel()
//	core.CurrentDashboardAdminPanel = a.DashboardAdminPanel
//	a.Config = core.NewConfig("configs/" + "test" + ".yml")
//	a.Config.InTests = true
//	core.CurrentConfig = a.Config
//	a.CommandRegistry = &CommandRegistry{
//		Actions: make(map[string]core.ICommand),
//	}
//	a.BlueprintRegistry = core.NewBlueprintRegistry()
//	a.Database = core.NewDatabase(a.Config)
//	a.Router = gin.Default()
//	a.Router.Use(cors.New(cors.Config{
//		AllowOrigins:     []string{"https://foo.com"},
//		AllowMethods:     []string{"PUT", "PATCH"},
//		AllowHeaders:     []string{"Origin"},
//		ExposeHeaders:    []string{"Content-Length"},
//		AllowCredentials: true,
//		AllowOriginFunc: func(origin string) bool {
//			return origin == "https://github.com"
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//	a.RegisterBaseCommands()
//	core.CurrentDatabaseSettings = &core.DatabaseSettings{
//		Default: a.Config.D.Db.Default,
//		Slave: a.Config.D.Db.Slave,
//	}
//	database := core.NewDatabaseInstance()
//	database.Adapter.SetIsolationLevelForTests(database.Db)
//	a.BlueprintRegistry.ResetMigrationTree()
//	StoreCurrentApp(&a)
//	if core.CurrentConfig.D.Db.Default.Type == "sqlite" {
//		appInstance.BlueprintRegistry.ResetMigrationTree()
//		downCommand := MigrateCommand{}
//		downCommand.Proceed("down", make([]string, 0))
//	}
//	upCommand := MigrateCommand{}
//	upCommand.Proceed("up", make([]string, 0))
//	return &a
//}

// Helper function to process a request and test its response
func TestHTTPResponse(t *testing.T, app *App, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	app.Router.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

var appForTests *App
var CreatedDatabaseForTests bool

func NewFullAppForTests() *App {
	if appForTests != nil {
		if appForTests.Config.D.Db.Default.Type == "sqlite" {
			appForTests.BlueprintRegistry.ResetMigrationTree()
			upCommand := MigrateCommand{}
			upCommand.Proceed("up", make([]string, 0))
		}
		return appForTests
	}
	a := NewApp(os.Getenv("TEST_ENVIRONMENT"), true)
	a.Config.InTests = true
	StoreCurrentApp(a)
	appForTests = a
	a.Initialize()
	authBlueprintInterface, _ := a.BlueprintRegistry.GetByName("auth")
	authBlueprint := authBlueprintInterface.(auth.Blueprint)
	authBlueprint.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.DirectAuthProvider{})
	authBlueprint.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.TokenAuthProvider{})
	authBlueprint.AuthAdapterRegistry.RegisterNewAdapter(&interfaces3.TokenWithExpirationAuthProvider{})
	a.InitializeRouter()
	// appForTests.DashboardAdminPanel.RegisterHTTPHandlers(a.Router)
	return a
}

//type GoMonolithTestSuite struct {
//	suite.Suite
//}
//
//func (suite *GoMonolithTestSuite) SetupTest() {
//	db := dialect.GetDB()
//	db = db.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
//	if db.Error != nil {
//		assert.Equal(suite.T(), true, false, "Couldnt setup isolation level for db")
//	}
//	db = db.Exec("BEGIN")
//	if db.Error != nil {
//		assert.Equal(suite.T(), true, false, "Couldnt start transaction")
//	}
//}
//
//func (suite *GoMonolithTestSuite) TearDownSuite() {
//	db := dialect.GetDB()
//	db = db.Exec("ROLLBACK")
//	if db.Error != nil {
//		assert.Equal(suite.T(), true, false, "Couldnt rollback transaction")
//	}
//}
