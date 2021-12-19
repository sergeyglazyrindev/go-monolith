package core

import (
	"github.com/gin-gonic/gin"
)

type IApp interface {
	GetConfig() *Config
	GetDatabase() *Database
	GetRouter() *gin.Engine
	GetCommandRegistry() *CommandRegistry
	GetBlueprintRegistry() IBlueprintRegistry
	GetDashboardAdminPanel() *DashboardAdminPanel
	GetAuthAdapterRegistry() *AuthProviderRegistry
}
