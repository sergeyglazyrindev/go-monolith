---
sidebar_position: 1
---

# Object query builder

Mostly it's used to filter records in admin panel.  
Each operator has to implement following interface:
```go
type IGormOperator interface {
	Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext
	GetName() string
	RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error
	TransformValue(value string) interface{}
}
```
You can register your own operator, using command:
```go
ProjectGormOperatorRegistry.RegisterOperator(&YOUROWNOPERATOR{})
```
