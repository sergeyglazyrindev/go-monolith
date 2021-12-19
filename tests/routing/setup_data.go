package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/sergeyglazyrindev/go-monolith/core"
)

type BlueprintRouting struct {
	core.Blueprint
}

var ConcreteBlueprint BlueprintRouting
var visited = false

func (b BlueprintRouting) InitRouter(app core.IApp, group *gin.RouterGroup) {
	// mainRouter *gin.Engine
	group.GET("/visit/", func(c *gin.Context) {
		visited = true
	})
}

func init() {
	ConcreteBlueprint = BlueprintRouting{
		core.Blueprint{
			Name:              "user",
			Description:       "blueprint",
			MigrationRegistry: core.NewMigrationRegistry(),
		},
	}
}
