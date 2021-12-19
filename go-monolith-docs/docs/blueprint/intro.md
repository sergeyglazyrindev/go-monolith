---
sidebar_position: 1
---

# Blueprint

Blueprint is an implementation of following interface
```go
type IBlueprint interface {
	GetName() string
	GetDescription() string
	GetMigrationRegistry() IMigrationRegistry
	InitRouter(app IApp, group *gin.RouterGroup)
	InitApp(app IApp)
}
```
There's a command to add blueprint to your project.  
Please initialize everything not related to http in the InitApp method that is called during app initialization.
Each blueprint would have its own gin RouterGroup, you may do add handlers to mainRouter as well.  
But try to avoid that.  
Also, don't forget to register blueprint in your app, like here
```go
a.BlueprintRegistry.Register(userblueprint.ConcreteBlueprint)
```
