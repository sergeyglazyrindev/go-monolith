package core

import (
	"gorm.io/gorm"
	"time"
)

type IUser interface {
	GetID() uint
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() gorm.DeletedAt
	GetUsername() string
	GetFirstName() string
	GetLastName() string
	GetPassword() string
	GetIsPasswordUsable() bool
	GetEmail() string
	GetActive() bool
	GetIsStaff() bool
	GetIsSuperUser() bool
	GetUserGroups() *[]UserGroup
	GetPermissions() *[]Permission
	GetPhoto() string
	GetLastLogin() *time.Time
	GetExpiresOn() *time.Time
	GetGeneratedOTPToVerify() string
	GetOTPSeed() string
	GetOTPRequired() bool
	GetSalt() string
	GetPermissionRegistry() *UserPermRegistry
	SetCreatedAt(t *time.Time)
	SetUpdatedAt(t *time.Time)
	SetDeletedAt(t gorm.DeletedAt)
	SetUsername(username string)
	SetFirstName(firstName string)
	SetLastName(lastName string)
	SetPassword(password string)
	SetIsPasswordUsable(isPasswordUsable bool)
	SetEmail(email string)
	SetActive(isActive bool)
	SetIsStaff(isStaff bool)
	SetIsSuperUser(isSuperUser bool)
	SetUserGroups(userGroups *[]UserGroup)
	SetPermissions(permissions *[]Permission)
	SetPhoto(photo string)
	SetLastLogin(t *time.Time)
	SetExpiresOn(t *time.Time)
	SetGeneratedOTPToVerify(generatedOtpToVerify string)
	SetOTPSeed(seed string)
	SetOTPRequired(isOtpRequired bool)
	SetSalt(salt string)
	SetPermissionRegistry(upr *UserPermRegistry)
	String() string
	GetFullName() string
	BuildPermissionRegistry() *UserPermRegistry
}

var GenerateUserModel = func() IUser {
	return &User{}
}

var GenerateBunchOfUserModels = func() interface{} {
	return &[]*User{}
}
