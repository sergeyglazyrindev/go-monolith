package migrations

import (
	"github.com/sergeyglazyrindev/go-monolith/core"
)

type initial1621667392 struct {
}

func (m initial1621667392) GetName() string {
	return "user.1621667393"
}

func (m initial1621667392) GetID() int64 {
	return 1621667392
}

func (m initial1621667392) Up(database *core.ProjectDatabase) error {
	appliedMigrations = append(appliedMigrations, m.GetName())
	return nil
}

func (m initial1621667392) Down(database *core.ProjectDatabase) error {
	appliedMigrations = core.Remove(appliedMigrations, m.GetName())
	return nil
}

func (m initial1621667392) Deps() []string {
	return make([]string, 0)
}

type migration21621667393 struct {
}

func (m migration21621667393) GetName() string {
	return "user.1621680132"
}

func (m migration21621667393) GetID() int64 {
	return 1621667393
}

func (m migration21621667393) Up(database *core.ProjectDatabase) error {
	appliedMigrations = append(appliedMigrations, m.GetName())
	return nil
}

func (m migration21621667393) Down(database *core.ProjectDatabase) error {
	appliedMigrations = core.Remove(appliedMigrations, m.GetName())
	return nil
}

func (m migration21621667393) Deps() []string {
	return []string{"user.1621667393"}
}

type initialtest11621667392 struct {
}

func (m initialtest11621667392) GetName() string {
	return "test1.1621667393"
}

func (m initialtest11621667392) GetID() int64 {
	return 1621667394
}

func (m initialtest11621667392) Up(database *core.ProjectDatabase) error {
	appliedMigrations = append(appliedMigrations, m.GetName())
	return nil
}

func (m initialtest11621667392) Down(database *core.ProjectDatabase) error {
	appliedMigrations = core.Remove(appliedMigrations, m.GetName())
	return nil
}

func (m initialtest11621667392) Deps() []string {
	return make([]string, 0)
}

type migration2test1621667393 struct {
}

func (m migration2test1621667393) GetName() string {
	return "test1.1621680132"
}

func (m migration2test1621667393) GetID() int64 {
	return 1621667395
}

func (m migration2test1621667393) Up(database *core.ProjectDatabase) error {
	appliedMigrations = append(appliedMigrations, m.GetName())
	return nil
}

func (m migration2test1621667393) Down(database *core.ProjectDatabase) error {
	appliedMigrations = core.Remove(appliedMigrations, m.GetName())
	return nil
}

func (m migration2test1621667393) Deps() []string {
	return []string{"test1.1621667393", "user.1621680132"}
}

type initialblueprintconflicts1621667392 struct {
}

func (m initialblueprintconflicts1621667392) GetName() string {
	return "user.1621680132"
}

func (m initialblueprintconflicts1621667392) GetID() int64 {
	return 1621667396
}

func (m initialblueprintconflicts1621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m initialblueprintconflicts1621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m initialblueprintconflicts1621667392) Deps() []string {
	return make([]string, 0)
}

type migration2blueprintconflicts1621667392 struct {
}

func (m migration2blueprintconflicts1621667392) GetName() string {
	return "user.16216801321"
}

func (m migration2blueprintconflicts1621667392) GetID() int64 {
	return 1621667397
}

func (m migration2blueprintconflicts1621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m migration2blueprintconflicts1621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m migration2blueprintconflicts1621667392) Deps() []string {
	return []string{"user.1621680132"}
}

type migration3blueprintconflicts1621667392 struct {
}

func (m migration3blueprintconflicts1621667392) GetName() string {
	return "user.16216801341"
}

func (m migration3blueprintconflicts1621667392) GetID() int64 {
	return 1621667398
}

func (m migration3blueprintconflicts1621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m migration3blueprintconflicts1621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m migration3blueprintconflicts1621667392) Deps() []string {
	return []string{"user.16216801321"}
}

type migration4blueprintconflicts1621667392 struct {
}

func (m migration4blueprintconflicts1621667392) GetName() string {
	return "user.16216801381"
}

func (m migration4blueprintconflicts1621667392) GetID() int64 {
	return 1621667399
}

func (m migration4blueprintconflicts1621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m migration4blueprintconflicts1621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m migration4blueprintconflicts1621667392) Deps() []string {
	return []string{"user.16216801321"}
}

type nodeps11621667392 struct {
}

func (m nodeps11621667392) GetName() string {
	return "test1.1621680132"
}

func (m nodeps11621667392) GetID() int64 {
	return 1621667492
}

func (m nodeps11621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m nodeps11621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m nodeps11621667392) Deps() []string {
	return make([]string, 0)
}

type nodeps21621667392 struct {
}

func (m nodeps21621667392) GetName() string {
	return "test1.16216801321"
}

func (m nodeps21621667392) GetID() int64 {
	return 1621667493
}

func (m nodeps21621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m nodeps21621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m nodeps21621667392) Deps() []string {
	return make([]string, 0)
}

type loopedmigration11621667392 struct {
}

func (m loopedmigration11621667392) GetName() string {
	return "user.1621680132"
}

func (m loopedmigration11621667392) GetID() int64 {
	return 1621667592
}

func (m loopedmigration11621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m loopedmigration11621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m loopedmigration11621667392) Deps() []string {
	return []string{"user1.16216801321"}
}

type loopedmigration21621667392 struct {
}

func (m loopedmigration21621667392) GetName() string {
	return "user1.16216801321"
}

func (m loopedmigration21621667392) GetID() int64 {
	return 1621667593
}

func (m loopedmigration21621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m loopedmigration21621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m loopedmigration21621667392) Deps() []string {
	return []string{"user.1621680132"}
}

type samenamemigration11621667392 struct {
}

func (m samenamemigration11621667392) GetName() string {
	return "user.1621680132"
}

func (m samenamemigration11621667392) GetID() int64 {
	return 1621667692
}

func (m samenamemigration11621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m samenamemigration11621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m samenamemigration11621667392) Deps() []string {
	return make([]string, 0)
}

type samenamemigration21621667392 struct {
}

func (m samenamemigration21621667392) GetName() string {
	return "user.1621680132"
}

func (m samenamemigration21621667392) GetID() int64 {
	return 1621667693
}

func (m samenamemigration21621667392) Up(database *core.ProjectDatabase) error {
	return nil
}

func (m samenamemigration21621667392) Down(database *core.ProjectDatabase) error {
	return nil
}

func (m samenamemigration21621667392) Deps() []string {
	return make([]string, 0)
}

var TestBlueprintMigrationRegistry *core.MigrationRegistry
var Test1BlueprintMigrationRegistry *core.MigrationRegistry
var TestBlueprint core.Blueprint

var Test1Blueprint core.Blueprint
var BlueprintWithConflictsMigrationRegistry *core.MigrationRegistry
var BlueprintWithConflicts core.Blueprint

var BlueprintWithNoMigrationsRegistry *core.MigrationRegistry
var BlueprintWithNoMigrations core.Blueprint

var BlueprintWithTwoSameDeps core.Blueprint

var Blueprint1WithLoopedMigrations core.Blueprint
var Blueprint2WithLoopedMigrations core.Blueprint

var Blueprint1WithSameMigrationNames core.Blueprint
var Blueprint2WithSameMigrationNames core.Blueprint

var appliedMigrations = make([]string, 0)

func init() {
	TestBlueprintMigrationRegistry = core.NewMigrationRegistry()
	TestBlueprintMigrationRegistry.AddMigration(initial1621667392{})
	TestBlueprintMigrationRegistry.AddMigration(migration21621667393{})

	Test1BlueprintMigrationRegistry = core.NewMigrationRegistry()
	Test1BlueprintMigrationRegistry.AddMigration(initialtest11621667392{})
	Test1BlueprintMigrationRegistry.AddMigration(migration2test1621667393{})

	TestBlueprint = core.Blueprint{
		Name:              "user",
		Description:       "this blueprint for testing",
		MigrationRegistry: TestBlueprintMigrationRegistry,
	}
	Test1Blueprint = core.Blueprint{
		Name:              "test1",
		Description:       "this test1 blueprint for testing",
		MigrationRegistry: Test1BlueprintMigrationRegistry,
	}

	BlueprintWithTwoSameDeps = core.Blueprint{
		Name:              "user",
		Description:       "this blueprint for testing",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	BlueprintWithTwoSameDeps.GetMigrationRegistry().AddMigration(nodeps11621667392{})
	BlueprintWithTwoSameDeps.GetMigrationRegistry().AddMigration(nodeps21621667392{})

	BlueprintWithConflictsMigrationRegistry = core.NewMigrationRegistry()
	BlueprintWithConflictsMigrationRegistry.AddMigration(initialblueprintconflicts1621667392{})
	BlueprintWithConflictsMigrationRegistry.AddMigration(migration2blueprintconflicts1621667392{})
	BlueprintWithConflictsMigrationRegistry.AddMigration(migration3blueprintconflicts1621667392{})
	BlueprintWithConflictsMigrationRegistry.AddMigration(migration4blueprintconflicts1621667392{})
	BlueprintWithConflicts = core.Blueprint{
		Name:              "user",
		Description:       "blueprint with conflicts",
		MigrationRegistry: BlueprintWithConflictsMigrationRegistry,
	}

	BlueprintWithNoMigrationsRegistry = core.NewMigrationRegistry()
	BlueprintWithNoMigrations = core.Blueprint{
		Name:              "user",
		Description:       "blueprint with no migrations",
		MigrationRegistry: BlueprintWithNoMigrationsRegistry,
	}

	Blueprint1WithLoopedMigrations = core.Blueprint{
		Name:              "user",
		Description:       "blueprint with looped migrations 1",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	Blueprint1WithLoopedMigrations.GetMigrationRegistry().AddMigration(loopedmigration11621667392{})
	Blueprint2WithLoopedMigrations = core.Blueprint{
		Name:              "user1",
		Description:       "blueprint with looped migrations 2",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	Blueprint2WithLoopedMigrations.GetMigrationRegistry().AddMigration(loopedmigration21621667392{})

	Blueprint1WithSameMigrationNames = core.Blueprint{
		Name:              "user",
		Description:       "blueprint with same migration names",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	Blueprint2WithSameMigrationNames = core.Blueprint{
		Name:              "user1",
		Description:       "blueprint with same migration names",
		MigrationRegistry: core.NewMigrationRegistry(),
	}
	Blueprint1WithSameMigrationNames.GetMigrationRegistry().AddMigration(samenamemigration11621667392{})
	Blueprint2WithSameMigrationNames.GetMigrationRegistry().AddMigration(samenamemigration21621667392{})
}
