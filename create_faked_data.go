package gomonolith

import (
	"fmt"
	abtestmodel "github.com/sergeyglazyrindev/go-monolith/blueprint/abtest/models"
	"github.com/sergeyglazyrindev/go-monolith/blueprint/approval/models"
	logmodel "github.com/sergeyglazyrindev/go-monolith/blueprint/logging/models"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"strconv"
	"time"
)

type CreateFakedDataCommand struct {
}

func (c CreateFakedDataCommand) Proceed(subaction string, args []string) error {
	database := core.NewDatabaseInstance()
	for i := range core.GenerateNumberSequence(1, 100) {
		userModel := &core.User{
			Email:     fmt.Sprintf("admin_%d@example.com", i),
			Username:  "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName:  "lastname_" + strconv.Itoa(i),
		}
		database.Db.Create(&userModel)
		oneTimeAction := &core.OneTimeAction{
			User:      *userModel,
			ExpiresOn: time.Now(),
			Code:      strconv.Itoa(i),
		}
		database.Db.Create(&oneTimeAction)
		session := &core.Session{
			User:      userModel,
			LoginTime: time.Now(),
			LastLogin: time.Now(),
		}
		database.Db.Create(&session)
	}
	var contentTypes []*core.ContentType
	database.Db.Find(&contentTypes)
	for _, contentType := range contentTypes {
		logModel := logmodel.Log{
			Username:      "admin",
			ContentTypeID: contentType.ID,
		}
		database.Db.Create(&logModel)
		approvalModel := models.Approval{
			ContentTypeID:       contentType.ID,
			ModelPK:             uint(1),
			ColumnName:          "Email",
			OldValue:            "admin@example.com",
			NewValue:            "admin1@example.com",
			NewValueDescription: "changing email",
			ChangedBy:           "superuser",
			ChangeDate:          time.Now(),
		}
		database.Db.Create(&approvalModel)
		abTestModel := abtestmodel.ABTest{
			Name:          "test_1",
			ContentTypeID: contentType.ID,
		}
		database.Db.Create(&abTestModel)
		for i := range core.GenerateNumberSequence(0, 100) {
			abTestValueModel := abtestmodel.ABTestValue{
				ABTest: abTestModel,
				Value:  strconv.Itoa(i),
			}
			database.Db.Create(&abTestValueModel)
		}
	}
	for i := range core.GenerateNumberSequence(1, 20) {
		groupModel := core.UserGroup{
			GroupName: fmt.Sprintf("Group name %d", i),
		}
		database.Db.Create(&groupModel)
	}
	database.Close()
	return nil
}

func (c CreateFakedDataCommand) GetHelpText() string {
	return "Create fake data for testing go-monolith"
}
