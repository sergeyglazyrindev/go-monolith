package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type IAuthProvider interface {
	GetUserFromRequest(c *gin.Context) IUser
	Signin(c *gin.Context)
	Logout(c *gin.Context)
	IsAuthenticated(c *gin.Context)
	GetSession(c *gin.Context) ISessionProvider
	GetName() string
	Signup(c *gin.Context)
}

type AuthProviderRegistry struct {
	registeredAdapters map[string]IAuthProvider
}

func (r *AuthProviderRegistry) RegisterNewAdapter(adapter IAuthProvider) {
	r.registeredAdapters[adapter.GetName()] = adapter
}

func (r *AuthProviderRegistry) GetAdapter(name string) (IAuthProvider, error) {
	adapter, ok := r.registeredAdapters[name]
	if ok {
		return adapter, nil
	}
	return nil, fmt.Errorf("adapter with name %s not found", name)
}

func (r *AuthProviderRegistry) Iterate() <-chan IAuthProvider {
	chnl := make(chan IAuthProvider)
	go func() {
		defer close(chnl)
		for _, authProvider := range r.registeredAdapters {
			chnl <- authProvider
		}
	}()
	return chnl
}

func NewAuthProviderRegistry() *AuthProviderRegistry {
	return &AuthProviderRegistry{
		registeredAdapters: make(map[string]IAuthProvider),
	}
}
