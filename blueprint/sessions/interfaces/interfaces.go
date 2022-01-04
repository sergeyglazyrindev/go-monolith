package interfaces

/*
	Package contains session providers provided by GoMonolith by default, currently it's only DBSessionProvider, but
	it's easy to add another provider that works with different data storages, etc
*/
import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
)

type SessionProviderRegistry struct {
	registeredSessionAdapters map[string]core.ISessionProvider
	defaultAdapter            string
}

func (r *SessionProviderRegistry) RegisterNewAdapter(adapter core.ISessionProvider, defaultAdapter bool) {
	r.registeredSessionAdapters[adapter.GetName()] = adapter
	if defaultAdapter {
		r.defaultAdapter = adapter.GetName()
	}
}

func (r *SessionProviderRegistry) GetAdapter(name string) (core.ISessionProvider, error) {
	adapter, ok := r.registeredSessionAdapters[name]
	if ok {
		return adapter, nil
	}
	return nil, fmt.Errorf("adapter with name %s not found", name)
}

func (r *SessionProviderRegistry) GetDefaultAdapter() (core.ISessionProvider, error) {
	adapter, ok := r.registeredSessionAdapters[r.defaultAdapter]
	if ok {
		return adapter, nil
	}
	return nil, fmt.Errorf("no default session adapter configured")
}

func NewSessionRegistry() *SessionProviderRegistry {
	return &SessionProviderRegistry{
		registeredSessionAdapters: make(map[string]core.ISessionProvider),
		defaultAdapter:            "",
	}
}

func NewSession() *core.Session {
	key := core.GenerateBase64(24)
	return &core.Session{
		Key:  key,
		Data: "{}",
	}
}
