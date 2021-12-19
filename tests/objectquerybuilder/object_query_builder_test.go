package objectquerybuilder

import (
	"github.com/sergeyglazyrindev/go-monolith"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"math"
	"strconv"
	"testing"
	"time"
)

type ObjectQueryBuilderTestSuite struct {
	gomonolith.TestSuite
	createdUser *core.User
}

func (suite *ObjectQueryBuilderTestSuite) SetupTest() {
	suite.TestSuite.SetupTest()
	db := suite.Database.Db
	u := core.User{Username: "dsadas", Email: "ffsdfsd@example.com"}
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()
	db.Create(&u)
	suite.createdUser = &u
}

func (suite *ObjectQueryBuilderTestSuite) TestExact() {
	operator := core.ExactGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsadas", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsadaS", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestIExact() {
	operator := core.IExactGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsadas", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsadas", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestContains() {
	operator := core.ContainsGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsa", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, uint(0))
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsad", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
}

func (suite *ObjectQueryBuilderTestSuite) TestIContains() {
	operator := core.IContainsGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsad", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsa", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIn() {
	operator := core.InGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, []uint{suite.createdUser.ID}, core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestGt() {
	operator := core.GtGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, -1, core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestGte() {
	operator := core.GteGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, suite.createdUser.ID, core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestLt() {
	operator := core.LtGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, suite.createdUser.ID+100, core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestLte() {
	operator := core.LteGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, suite.createdUser.ID, core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestStartsWith() {
	operator := core.StartsWithGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsad", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsa", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestIStartsWith() {
	operator := core.IStartsWithGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "dsad", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Dsa", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestEndsWith() {
	operator := core.EndsWithGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "das", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Das", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestIEndsWith() {
	operator := core.IEndsWithGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "das", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]}, "Das", core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestRange() {
	operator := core.RangeGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["ID"]}, []uint{suite.createdUser.ID - 1, suite.createdUser.ID + 100}, core.NewSQLConditionBuilder("and"))
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestDate() {
	operator := core.DateGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]}, suite.createdUser.CreatedAt.Round(0).Format(time.RFC3339), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Round(0).Add(-10*3600*24*time.Second), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, uint(0))
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["UpdatedAt"]},
		suite.createdUser.CreatedAt.Round(0).Add(-10*3600*24*time.Second), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, uint(0))
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["UpdatedAt"]},
		suite.createdUser.CreatedAt.Round(0).Format(time.RFC3339), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestYear() {
	operator := core.YearGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Year(), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestMonth() {
	operator := core.MonthGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Month(), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestDay() {
	operator := core.DayGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Day(), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestWeek() {
	operator := core.WeekGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	_, isoWeek := suite.createdUser.CreatedAt.ISOWeek()
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		isoWeek, core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestWeekDay() {
	operator := core.WeekDayGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Weekday(), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestQuarter() {
	operator := core.QuarterGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		int(math.Ceil(float64(suite.createdUser.CreatedAt.Month()/3))), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestTime() {
	operator := core.TimeGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Format("15:04"), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestHour() {
	operator := core.HourGormOperator{}
	var u core.User
	var u1 core.User
	suite.Database.Db.First(&u1)
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Hour(), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestRegex() {
	operator := core.RegexGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]},
		"^dsadas$", core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]},
		"^dsaDAS1111111111$", core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, uint(0))
}

func (suite *ObjectQueryBuilderTestSuite) TestIRegex() {
	operator := core.IRegexGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]},
		"^dsadas$", core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
	u = core.User{}
	gormOperatorContext = core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["Username"]},
		"^dsaDAS$", core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestMinute() {
	operator := core.MinuteGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Minute(), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestSecond() {
	operator := core.SecondGormOperator{}
	var u core.User
	var u1 core.User
	suite.Database.Db.First(&u1)
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		suite.createdUser.CreatedAt.Second(), core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestIsNull() {
	operator := core.IsNullGormOperator{}
	var u core.User
	gormOperatorContext := core.NewGormOperatorContext(core.NewGormPersistenceStorage(suite.Database.Db), &core.User{})
	operator.Build(
		suite.Database.Adapter, gormOperatorContext, &core.Field{Field: *gormOperatorContext.Statement.Schema.FieldsByName["CreatedAt"]},
		false, core.NewSQLConditionBuilder("and"),
	)
	gormOperatorContext.Tx.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestFilterGormModel() {
	var u core.User
	createdSecond := suite.createdUser.CreatedAt.Second()
	statement := &gorm.Statement{DB: suite.Database.Db}
	statement.Parse(&core.User{})
	core.FilterGormModel(suite.Database.Adapter, core.NewGormPersistenceStorage(suite.Database.Db), statement.Schema, []string{"CreatedAt__isnull=false", "CreatedAt__second=" + strconv.Itoa(createdSecond)}, &core.User{})
	suite.Database.Db.First(&u)
	assert.Equal(suite.T(), u.ID, suite.createdUser.ID)
}

func (suite *ObjectQueryBuilderTestSuite) TestFilterGormModelWithDependencies() {
	var p core.Permission
	statement := &gorm.Statement{DB: suite.Database.Db}
	statement.Parse(&core.Permission{})
	contentType := core.ContentType{BlueprintName: "user11", ModelName: "user11"}
	suite.Database.Db.Create(&contentType)
	permission := core.Permission{ContentType: contentType, PermissionBits: core.RevertPermBit}
	suite.Database.Db.Create(&permission)
	core.FilterGormModel(suite.Database.Adapter, core.NewGormPersistenceStorage(suite.Database.Db), statement.Schema, []string{"CreatedAt__isnull=false", "ContentType__ID__exact=" + strconv.Itoa(int(contentType.ID))}, &core.Permission{})
	suite.Database.Db.First(&p)
	assert.Equal(suite.T(), p.ContentTypeID, contentType.ID)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestObjectQueryBuilder(t *testing.T) {
	gomonolith.RunTests(t, new(ObjectQueryBuilderTestSuite))
}
