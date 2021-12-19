package core

import (
	"fmt"
	"strings"
)

type PermissionDescribed struct {
	Name CustomPermission
	Bit  PermBitInteger
}

type IPermissionRegistry interface {
	AddPermission(permission CustomPermission, permissionBit PermBitInteger)
	GetPermissionBit(permission CustomPermission) PermBitInteger
	GetAllPermissions() <-chan *PermissionDescribed
}

type IUserPermissionRegistry interface {
	DoesUserHaveRightFor(permissionName string) bool
}

type CustomPermission string

type UserPermRegistry struct {
	BlueprintPerm map[string]*UserPerm
	IsSuperUser   bool
}

func getPermissionAliasForBlueprintAndModel(blueprintName string, modelName string) string {
	return fmt.Sprintf("b.%s.m.%s", blueprintName, modelName)
}

func (upr *UserPermRegistry) AddPermissionForBlueprint(blueprintName string, modelName string, userPerm *UserPerm) {
	blueprintModelAlias := getPermissionAliasForBlueprintAndModel(blueprintName, modelName)
	upr.BlueprintPerm[blueprintModelAlias] = userPerm
}

func (upr *UserPermRegistry) GetPermissionForBlueprint(blueprintName string, modelName string) *UserPerm {
	blueprintModelAlias := getPermissionAliasForBlueprintAndModel(blueprintName, modelName)
	userPerm, isExists := upr.BlueprintPerm[blueprintModelAlias]
	if !isExists {
		userPerm := NewUserPerm(0)
		userPerm.IsSuperUser = upr.IsSuperUser
		upr.AddPermissionForBlueprint(blueprintName, modelName, userPerm)
		return userPerm
	}
	return userPerm
}

func (upr *UserPermRegistry) IsThereAnyPermissionForBlueprint(blueprintName string) bool {
	if upr.IsSuperUser {
		return true
	}
	for userPermIdentifier := range upr.BlueprintPerm {
		if strings.HasSuffix(userPermIdentifier, fmt.Sprintf("b.%s.", blueprintName)) {
			return true
		}
	}
	return false
}

type PermRegistry struct {
	PermNameBitInteger map[CustomPermission]PermBitInteger
}

func (ap *PermRegistry) AddPermission(permission CustomPermission, permissionBit PermBitInteger) {
	_, isRegistered := ap.PermNameBitInteger[permission]
	if isRegistered {
		Trail(WARNING, "you are overriding permission with name %s", permission)
	}
	ap.PermNameBitInteger[permission] = permissionBit
}

func (ap *PermRegistry) GetPermissionName(permissionBit PermBitInteger) CustomPermission {
	for permissionName, permissionBitTmp := range ap.PermNameBitInteger {
		if permissionBitTmp == permissionBit {
			return permissionName
		}
	}
	return ""
}

func (ap *PermRegistry) GetPermissionBit(permission CustomPermission) PermBitInteger {
	permissionBit, isRegistered := ap.PermNameBitInteger[permission]
	if !isRegistered {
		Trail(CRITICAL, "no permission registered with name %s", permission)
		panic(fmt.Errorf("no permission registered with name %s", permission))
	}
	return permissionBit
}

func (ap *PermRegistry) GetAllPermissions() <-chan *PermissionDescribed {
	chnl := make(chan *PermissionDescribed)
	go func() {
		defer close(chnl)
		for permName, permBit := range ap.PermNameBitInteger {
			chnl <- &PermissionDescribed{Name: permName, Bit: permBit}
		}
	}()
	return chnl
}

type UserPerm struct {
	PermBitInteger PermBitInteger
	IsSuperUser    bool
}

func (ap *UserPerm) DoesUserHaveRightFor(permissionName CustomPermission) bool {
	if ap.IsSuperUser {
		return true
	}
	permissionBit := ProjectPermRegistry.GetPermissionBit(permissionName)
	return (ap.PermBitInteger & permissionBit) == permissionBit
}

func (ap *UserPerm) AddPermission(permBitInteger PermBitInteger) {
	ap.PermBitInteger = ap.PermBitInteger | permBitInteger
}

func (ap *UserPerm) HasReadPermission() bool {
	if ap.IsSuperUser {
		return true
	}
	return (ap.PermBitInteger & ReadPermBit) == ReadPermBit
}

func (ap *UserPerm) HasAddPermission() bool {
	if ap.IsSuperUser {
		return true
	}
	return (ap.PermBitInteger & AddPermBit) == AddPermBit
}

func (ap *UserPerm) HasEditPermission() bool {
	if ap.IsSuperUser {
		return true
	}
	return (ap.PermBitInteger & EditPermBit) == EditPermBit
}

func (ap *UserPerm) HasDeletePermission() bool {
	if ap.IsSuperUser {
		return true
	}
	return (ap.PermBitInteger & DeletePermBit) == DeletePermBit
}

func (ap *UserPerm) HasPublishPermission() bool {
	if ap.IsSuperUser {
		return true
	}
	return (ap.PermBitInteger & PublishPermBit) == PublishPermBit
}

func (ap *UserPerm) HasRevertPermission() bool {
	if ap.IsSuperUser {
		return true
	}
	return (ap.PermBitInteger & RevertPermBit) == RevertPermBit
}

func NewPerm() *PermRegistry {
	return &PermRegistry{
		PermNameBitInteger: make(map[CustomPermission]PermBitInteger),
	}
}

type PermBitInteger uint

const ReadPermBit PermBitInteger = 0
const AddPermBit PermBitInteger = 2
const EditPermBit PermBitInteger = 4
const DeletePermBit PermBitInteger = 8
const PublishPermBit PermBitInteger = 16
const RevertPermBit PermBitInteger = 32

var ProjectPermRegistry *PermRegistry

func init() {
	ProjectPermRegistry = NewPerm()
	ProjectPermRegistry.AddPermission("read", ReadPermBit)
	ProjectPermRegistry.AddPermission("add", AddPermBit)
	ProjectPermRegistry.AddPermission("edit", EditPermBit)
	ProjectPermRegistry.AddPermission("delete", DeletePermBit)
	ProjectPermRegistry.AddPermission("publish", PublishPermBit)
	ProjectPermRegistry.AddPermission("revert", RevertPermBit)
}

func NewUserPerm(permBitInteger PermBitInteger) *UserPerm {
	return &UserPerm{PermBitInteger: permBitInteger}
}

func NewUserPermRegistry() *UserPermRegistry {
	return &UserPermRegistry{BlueprintPerm: make(map[string]*UserPerm)}
}
