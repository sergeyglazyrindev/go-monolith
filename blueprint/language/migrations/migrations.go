package migrations

import (
	"github.com/sergeyglazyrindev/go-monolith/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
	BMigrationRegistry = core.NewMigrationRegistry()

	BMigrationRegistry.AddMigration(initial1623083053{})

	BMigrationRegistry.AddMigration(createall1623263607{})
	// placeholder to insert next migration
}
