package migrations

import (
	"errors"
	mapset "github.com/deckarep/golang-set"
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

type MigrationTestSuite struct {
	gomonolith.TestSuite
}

func (suite *MigrationTestSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

func (suite *MigrationTestSuite) BeforeTest(suiteName string, testName string) {
	appliedMigrations = make([]string, 0)
	suite.Database.Db.Exec("DROP TABLE migrations")
	suite.Database.Db.AutoMigrate(gomonolith.Migration{})
	suite.App.BlueprintRegistry = core.NewBlueprintRegistry()
	suite.App.BlueprintRegistry.Register(TestBlueprint)
	suite.App.BlueprintRegistry.Register(Test1Blueprint)
}

func (suite *MigrationTestSuite) TearDownSuite() {
	//err := os.Remove(suite.app.Config.D.Db.Default.Name)
	//if err != nil {
	//	assert.Equal(suite.T(), true, false, fmt.Errorf("Couldnt remove db with name %s", suite.app.Config.D.Db.Default.Name))
	//}
	gomonolith.ClearTestApp()
}

func (suite *MigrationTestSuite) TestUpgradeDatabase() {
	suite.App.BlueprintRegistry.Register(TestBlueprint)
	suite.App.BlueprintRegistry.Register(Test1Blueprint)
	suite.App.TriggerCommandExecution("migrate", "up", make([]string, 0))
	appliedMigrationsExpected := mapset.NewSet()
	appliedMigrationsExpected.Add("user.1621667393")
	appliedMigrationsExpected.Add("user.1621680132")
	appliedMigrationsExpected.Add("test1.1621667393")
	appliedMigrationsExpected.Add("test1.1621680132")
	appliedMigrationsActual := mapset.NewSet()
	for _, migrationName := range appliedMigrations {
		appliedMigrationsActual.Add(migrationName)
	}
	assert.Equal(suite.T(), appliedMigrationsExpected, appliedMigrationsActual)
	var appliedMigrationsDb []gomonolith.Migration
	suite.Database.Db.Find(&appliedMigrationsDb)
	assert.Equal(suite.T(), 4, len(appliedMigrationsDb))
}

func (suite *MigrationTestSuite) TestDowngradeDatabase() {
	suite.App.BlueprintRegistry.Register(TestBlueprint)
	suite.App.BlueprintRegistry.Register(Test1Blueprint)
	appliedMigrationsExpected := mapset.NewSet()
	appliedMigrationsExpected.Add("user.1621667393")
	appliedMigrationsExpected.Add("user.1621680132")
	appliedMigrationsExpected.Add("test1.1621667393")
	appliedMigrationsExpected.Add("test1.1621680132")
	suite.Database.Db.Create(
		&gomonolith.Migration{MigrationName: "user.1621667393", AppliedAt: time.Now()},
	)
	suite.Database.Db.Create(
		&gomonolith.Migration{MigrationName: "user.1621680132", AppliedAt: time.Now()},
	)
	suite.Database.Db.Create(
		&gomonolith.Migration{MigrationName: "test1.1621667393", AppliedAt: time.Now()},
	)
	suite.Database.Db.Create(
		&gomonolith.Migration{MigrationName: "test1.1621680132", AppliedAt: time.Now()},
	)
	var appliedMigrationsDb []gomonolith.Migration
	suite.Database.Db.Find(&appliedMigrationsDb)
	assert.Equal(suite.T(), 4, len(appliedMigrationsDb))
	suite.App.TriggerCommandExecution("migrate", "down", []string{""})
	appliedMigrationsDb = make([]gomonolith.Migration, 0)
	suite.Database.Db.Find(&appliedMigrationsDb)
	assert.Equal(suite.T(), 0, len(appliedMigrationsDb))
}

func (suite *MigrationTestSuite) TestTraverseDatabaseForUpgrade() {
	concreteBlueprintRegistry := core.NewBlueprintRegistry()
	concreteBlueprintRegistry.Register(TestBlueprint)
	concreteBlueprintRegistry.Register(Test1Blueprint)
	for res := range concreteBlueprintRegistry.TraverseMigrations() {
		res.Node.Apply(suite.Database)
	}
	appliedMigrationsExpected := mapset.NewSet()
	appliedMigrationsExpected.Add("user.1621667393")
	appliedMigrationsExpected.Add("user.1621680132")
	appliedMigrationsExpected.Add("test1.1621667393")
	appliedMigrationsExpected.Add("test1.1621680132")
	appliedMigrationsActual := mapset.NewSet()
	for _, migrationName := range appliedMigrations {
		appliedMigrationsActual.Add(migrationName)
	}
	assert.Equal(suite.T(), appliedMigrationsExpected, appliedMigrationsActual)
}

func (suite *MigrationTestSuite) TestTraverseDatabaseForDowngrade() {
	concreteBlueprintRegistry := core.NewBlueprintRegistry()
	concreteBlueprintRegistry.Register(TestBlueprint)
	concreteBlueprintRegistry.Register(Test1Blueprint)
	toDowngradeMigrationsExpected := mapset.NewSet()
	toDowngradeMigrationsExpected.Add("user.1621667393")
	toDowngradeMigrationsExpected.Add("user.1621680132")
	toDowngradeMigrationsExpected.Add("test1.1621667393")
	toDowngradeMigrationsExpected.Add("test1.1621680132")
	downgradedMigrationsActual := mapset.NewSet()
	for res := range concreteBlueprintRegistry.TraverseMigrationsDownTo(0) {
		res.Node.Downgrade(suite.Database)
		downgradedMigrationsActual.Add(res.Node.GetMigration().GetName())
	}
	assert.Equal(suite.T(), toDowngradeMigrationsExpected, downgradedMigrationsActual)
}

func (suite *MigrationTestSuite) TestBuildTreeForBlueprintWithNoMigrations() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(BlueprintWithNoMigrations)
	for res := range blueprintRegistry.TraverseMigrations() {
		if res.Error != nil {
			assert.Equal(suite.T(), false, true, res.Error)
			return
		}
	}
}

func (suite *MigrationTestSuite) TestBuildTreeWithNoUserBlueprint() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(Test1Blueprint)
	for res := range blueprintRegistry.TraverseMigrations() {
		assert.Equal(suite.T(), res.Error, errors.New("Couldn't find blueprint with name user"))
		return
	}
	assert.True(suite.T(), false)
}

func (suite *MigrationTestSuite) TestBuildTreeWithTwoNoDepsMigrationsFromtheSameBlueprint() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(BlueprintWithTwoSameDeps)
	for res := range blueprintRegistry.TraverseMigrations() {
		assert.True(suite.T(), strings.Contains(res.Error.Error(), "Found two or more migrations with no children from the same blueprint"))
		return
	}
	assert.True(suite.T(), false)
}

func (suite *MigrationTestSuite) TestBuildTreeWithTwoNoChildMigrationsFromtheSameBlueprint() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(BlueprintWithConflicts)
	for res := range blueprintRegistry.TraverseMigrations() {
		assert.True(suite.T(), strings.Contains(res.Error.Error(), "Found two or more migrations with no children from the same blueprint"))
		return
	}
	assert.True(suite.T(), false)
}

func (suite *MigrationTestSuite) TestBuildTreeWithLoop() {
	blueprintRegistry := core.NewBlueprintRegistry()
	blueprintRegistry.Register(Blueprint1WithLoopedMigrations)
	blueprintRegistry.Register(Blueprint2WithLoopedMigrations)
	for range blueprintRegistry.TraverseMigrations() {
	}
}

//func (suite *MigrationTestSuite) TestBuildTreeWithTwoSameMigrationNames() {
//	blueprintRegistry := core.NewBlueprintRegistry()
//	blueprintRegistry.Register(Blueprint1WithSameMigrationNames)
//	blueprintRegistry.Register(Blueprint2WithSameMigrationNames)
//	for res := range blueprintRegistry.TraverseMigrations() {
//		assert.True(suite.T(), strings.Contains(res.Error.Error(), "has been added to tree before"))
//	}
//}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSqliteMigrations(t *testing.T) {
	gomonolith.ClearApp()
	gomonolith.RunTests(t, new(MigrationTestSuite))
}
