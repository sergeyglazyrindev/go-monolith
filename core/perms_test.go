package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPerm(t *testing.T) {
	perm := NewUserPerm(ReadPermBit | AddPermBit | EditPermBit | DeletePermBit | PublishPermBit | RevertPermBit)
	assert.True(t, perm.HasAddPermission())
	assert.True(t, perm.HasReadPermission())
	assert.True(t, perm.HasEditPermission())
	assert.True(t, perm.HasDeletePermission())
	assert.True(t, perm.HasPublishPermission())
	assert.True(t, perm.HasRevertPermission())
	assert.True(t, perm.DoesUserHaveRightFor("read"))
}

func TestUserPermRegistry(t *testing.T) {
	userPermRegistry := NewUserPermRegistry()
	perm := NewUserPerm(ReadPermBit | AddPermBit | EditPermBit | DeletePermBit | PublishPermBit | RevertPermBit)
	userPermRegistry.AddPermissionForBlueprint("user", "user", perm)
	perm = userPermRegistry.GetPermissionForBlueprint("user", "user")
	assert.True(t, perm.HasAddPermission())
	assert.True(t, perm.HasReadPermission())
	assert.True(t, perm.HasEditPermission())
	assert.True(t, perm.HasDeletePermission())
	assert.True(t, perm.HasPublishPermission())
	assert.True(t, perm.HasRevertPermission())
	assert.True(t, perm.DoesUserHaveRightFor("read"))
}
