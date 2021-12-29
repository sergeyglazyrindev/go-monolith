package core

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strconv"
	"strings"
)

type GormOperatorContext struct {
	TableName string
	Tx        IPersistenceStorage
	Statement *gorm.Statement
}

// tx *gorm.DB, field *schema.Field

type IRegisterDbHandler interface {
	RegisterFunc(name string, impl interface{}, pure bool) error
}

type IGormOperator interface {
	Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext
	GetName() string
	RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error
	TransformValue(value string) interface{}
}

type ExactGormOperator struct {
}

func (ego *ExactGormOperator) GetName() string {
	return "exact"
}

func (ego *ExactGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *ExactGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *ExactGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	value1 := TransformValueForOperator(value)
	adapter.Exact(context, field, value1, SQLConditionBuilder)
	return context
}

type ArrayIncludesGormOperator struct {

}

func (ego *ArrayIncludesGormOperator) GetName() string {
	return "arrayin"
}

func (ego *ArrayIncludesGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *ArrayIncludesGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *ArrayIncludesGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	value1 := TransformValueForOperator(value)
	adapter.ArrayIncludes(context, field, value1, SQLConditionBuilder)
	return context
}

type JSONIncludesGormOperator struct {

}

func (ego *JSONIncludesGormOperator) GetName() string {
	return "jsoncontains"
}

func (ego *JSONIncludesGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *JSONIncludesGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *JSONIncludesGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	value1 := TransformValueForOperator(value)
	adapter.JSONContains(context, field, value1, SQLConditionBuilder)
	return context
}

type IExactGormOperator struct {
}

func (ego *IExactGormOperator) GetName() string {
	return "iexact"
}

func (ego *IExactGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *IExactGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *IExactGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	value1 := TransformValueForOperator(value)
	adapter.IExact(context, field, value1, SQLConditionBuilder)
	return context
}

type ContainsGormOperator struct {
}

func (ego *ContainsGormOperator) GetName() string {
	return "contains"
}

func (ego *ContainsGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *ContainsGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *ContainsGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Contains(context, field, value, SQLConditionBuilder)
	return context
}

type IContainsGormOperator struct {
}

func (ego *IContainsGormOperator) GetName() string {
	return "icontains"
}

func (ego *IContainsGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *IContainsGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *IContainsGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.IContains(context, field, value, SQLConditionBuilder)
	return context
}

type InGormOperator struct {
}

func (ego *InGormOperator) GetName() string {
	return "in"
}

func (ego *InGormOperator) TransformValue(value string) interface{} {
	return strings.Split(value, ",")
}

func (ego *InGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *InGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.In(context, field, value, SQLConditionBuilder)
	return context
}

type GtGormOperator struct {
}

func (ego *GtGormOperator) GetName() string {
	return "gt"
}

func (ego *GtGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *GtGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *GtGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Gt(context, field, value, SQLConditionBuilder)
	return context
}

type GteGormOperator struct {
}

func (ego *GteGormOperator) GetName() string {
	return "gte"
}

func (ego *GteGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *GteGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *GteGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Gte(context, field, value, SQLConditionBuilder)
	return context
}

type LtGormOperator struct {
}

func (ego *LtGormOperator) GetName() string {
	return "lt"
}

func (ego *LtGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *LtGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *LtGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Lt(context, field, value, SQLConditionBuilder)
	return context
}

type LteGormOperator struct {
}

func (ego *LteGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *LteGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *LteGormOperator) GetName() string {
	return "lte"
}

func (ego *LteGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Lte(context, field, value, SQLConditionBuilder)
	return context
}

type StartsWithGormOperator struct {
}

func (ego *StartsWithGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *StartsWithGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *StartsWithGormOperator) GetName() string {
	return "startswith"
}

func (ego *StartsWithGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.StartsWith(context, field, value, SQLConditionBuilder)
	return context
}

type IStartsWithGormOperator struct {
}

func (ego *IStartsWithGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *IStartsWithGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *IStartsWithGormOperator) GetName() string {
	return "istartswith"
}

func (ego *IStartsWithGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.IStartsWith(context, field, value, SQLConditionBuilder)
	return context
}

type EndsWithGormOperator struct {
}

func (ego *EndsWithGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *EndsWithGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *EndsWithGormOperator) GetName() string {
	return "endswith"
}

func (ego *EndsWithGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.EndsWith(context, field, value, SQLConditionBuilder)
	return context
}

type IEndsWithGormOperator struct {
}

func (ego *IEndsWithGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *IEndsWithGormOperator) GetName() string {
	return "iendswith"
}

func (ego *IEndsWithGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *IEndsWithGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.IEndsWith(context, field, value, SQLConditionBuilder)
	return context
}

type RangeGormOperator struct {
}

func (ego *RangeGormOperator) TransformValue(value string) interface{} {
	return strings.Split(value, ",")
}

func (ego *RangeGormOperator) GetName() string {
	return "range"
}

func (ego *RangeGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *RangeGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Range(context, field, value, SQLConditionBuilder)
	return context
}

type DateGormOperator struct {
}

func (ego *DateGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *DateGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *DateGormOperator) GetName() string {
	return "date"
}

func (ego *DateGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Date(context, field, value, SQLConditionBuilder)
	return context
}

type YearGormOperator struct {
}

func (ego *YearGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *YearGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *YearGormOperator) GetName() string {
	return "year"
}

func (ego *YearGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Year(context, field, value, SQLConditionBuilder)
	return context
}

type MonthGormOperator struct {
}

func (ego *MonthGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *MonthGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *MonthGormOperator) GetName() string {
	return "month"
}

func (ego *MonthGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Month(context, field, value, SQLConditionBuilder)
	return context
}

type DayGormOperator struct {
}

func (ego *DayGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *DayGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *DayGormOperator) GetName() string {
	return "day"
}

func (ego *DayGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Day(context, field, value, SQLConditionBuilder)
	return context
}

type WeekGormOperator struct {
}

func (ego *WeekGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *WeekGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *WeekGormOperator) GetName() string {
	return "week"
}

func (ego *WeekGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Week(context, field, value, SQLConditionBuilder)
	return context
}

type WeekDayGormOperator struct {
}

func (ego *WeekDayGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *WeekDayGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *WeekDayGormOperator) GetName() string {
	return "week_day"
}

func (ego *WeekDayGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.WeekDay(context, field, value, SQLConditionBuilder)
	return context
}

type QuarterGormOperator struct {
}

func (ego *QuarterGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *QuarterGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *QuarterGormOperator) GetName() string {
	return "quarter"
}

func (ego *QuarterGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Quarter(context, field, value, SQLConditionBuilder)
	return context
}

type TimeGormOperator struct {
}

func (ego *TimeGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *TimeGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *TimeGormOperator) GetName() string {
	return "time"
}

func (ego *TimeGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Time(context, field, value, SQLConditionBuilder)
	return context
}

type HourGormOperator struct {
}

func (ego *HourGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *HourGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *HourGormOperator) GetName() string {
	return "hour"
}

func (ego *HourGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Hour(context, field, value, SQLConditionBuilder)
	return context
}

type MinuteGormOperator struct {
}

func (ego *MinuteGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *MinuteGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *MinuteGormOperator) GetName() string {
	return "minute"
}

func (ego *MinuteGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Minute(context, field, value, SQLConditionBuilder)
	return context
}

type SecondGormOperator struct {
}

func (ego *SecondGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *SecondGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *SecondGormOperator) GetName() string {
	return "second"
}

func (ego *SecondGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Second(context, field, value, SQLConditionBuilder)
	return context
}

type IsNullGormOperator struct {
}

func (ego *IsNullGormOperator) TransformValue(value string) interface{} {
	isTruthy, _ := strconv.ParseBool(value)
	return isTruthy
}

func (ego *IsNullGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *IsNullGormOperator) GetName() string {
	return "isnull"
}

func (ego *IsNullGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.IsNull(context, field, value, SQLConditionBuilder)
	return context
}

type RegexGormOperator struct {
}

func (ego *RegexGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *RegexGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *RegexGormOperator) GetName() string {
	return "regex"
}

func (ego *RegexGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.Regex(context, field, value, SQLConditionBuilder)
	return context
}

type IRegexGormOperator struct {
}

func (ego *IRegexGormOperator) TransformValue(value string) interface{} {
	return value
}

func (ego *IRegexGormOperator) RegisterDbHandlers(registerDbHandler IRegisterDbHandler) error {
	return nil
}

func (ego *IRegexGormOperator) GetName() string {
	return "iregex"
}

func (ego *IRegexGormOperator) Build(adapter IDbAdapter, context *GormOperatorContext, field *Field, value interface{}, SQLConditionBuilder ISQLConditionBuilder) *GormOperatorContext {
	adapter.IRegex(context, field, value, SQLConditionBuilder)
	return context
}

type GormOperatorRegistry struct {
	Operators map[string]IGormOperator
}

func (gor *GormOperatorRegistry) RegisterOperator(operator IGormOperator) {
	gor.Operators[operator.GetName()] = operator
}

func (gor *GormOperatorRegistry) GetOperatorByName(operatorName string) (IGormOperator, error) {
	operator, exists := gor.Operators[strings.ToLower(operatorName)]
	if !exists {
		return nil, fmt.Errorf("no operator with name %s registered", strings.ToLower(operatorName))
	}
	return operator, nil
}

func (gor *GormOperatorRegistry) GetAll() <-chan IGormOperator {
	chnl := make(chan IGormOperator)
	go func() {
		defer close(chnl)
		for _, operator := range gor.Operators {
			chnl <- operator
		}
	}()
	return chnl
}

type GormQueryBuilder struct {
	FilterString     []string
	OperatorRegistry *GormOperatorRegistry
	Model            interface{}
}

var ProjectGormOperatorRegistry *GormOperatorRegistry

func init() {
	ProjectGormOperatorRegistry = &GormOperatorRegistry{
		Operators: make(map[string]IGormOperator),
	}
	ProjectGormOperatorRegistry.RegisterOperator(&ExactGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&IExactGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&JSONIncludesGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&ArrayIncludesGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&ContainsGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&IContainsGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&InGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&GtGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&GteGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&LtGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&LteGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&StartsWithGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&IStartsWithGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&EndsWithGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&IEndsWithGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&RangeGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&DateGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&YearGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&MonthGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&DayGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&WeekGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&WeekDayGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&QuarterGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&TimeGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&HourGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&MinuteGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&SecondGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&IsNullGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&RegexGormOperator{})
	ProjectGormOperatorRegistry.RegisterOperator(&IRegexGormOperator{})
}

func NewGormOperatorContext(db IPersistenceStorage, model interface{}) *GormOperatorContext {
	statement := &gorm.Statement{DB: db.(*GormPersistenceStorage).Db}
	statement.Parse(model)
	return &GormOperatorContext{
		Tx:        db,
		TableName: statement.Table,
		Statement: statement,
	}
}

func FilterGormModel(adapter IDbAdapter, db IPersistenceStorage, schema1 *schema.Schema, filterString []string, model interface{}) *GormOperatorContext {
	context := NewGormOperatorContext(db, model)
	context.Tx = db
	gormModelV := reflect.Indirect(reflect.ValueOf(model))
	for _, filter := range filterString {
		filterParams := strings.Split(filter, "=")
		filterName := filterParams[0]
		filterValue := filterParams[1]
		filterNameParams := strings.Split(filterName, "__")
		field, _ := schema1.FieldsByName[filterNameParams[0]]
		field1 := NewGoMonolithFieldFromGormField(gormModelV, field, nil, false)
		if field.DBName == "" {
			joinModelI := ProjectModels.GetModelByName(field.FieldType.Name())
			relation := context.Statement.Schema.Relationships
			relationsString := []string{}
			for _, relation1 := range relation.Relations {
				for _, reference := range relation1.References {
					relationsString = append(
						relationsString,
						fmt.Sprintf(
							"%s.%s = %s.%s",
							joinModelI.Statement.Table, reference.PrimaryKey.DBName, context.Statement.Table,
							reference.ForeignKey.DBName,
						),
					)
				}
			}
			if field.NotNull {
				context.Tx = context.Tx.Joins(
					fmt.Sprintf(
						"INNER JOIN %s on %s",
						joinModelI.Statement.Table, strings.Join(relationsString, " AND "),
					),
				)
			} else {
				context.Tx = context.Tx.Joins(
					fmt.Sprintf(
						"LEFT JOIN %s on %s",
						joinModelI.Statement.Table, strings.Join(relationsString, " AND "),
					),
				)
			}
			filterRelation := strings.Replace(filter, field.Name+"__", "", 1)
			FilterGormModel(adapter, context.Tx, field.Schema, []string{filterRelation}, joinModelI.Model)
			continue
		}
		operator, _ := ProjectGormOperatorRegistry.GetOperatorByName(filterNameParams[len(filterNameParams)-1])
		filterValueTransformed := operator.TransformValue(filterValue)
		context = operator.Build(adapter, context, field1, filterValueTransformed, &SQLConditionBuilder{Type: "and"})
	}
	return context
}

type ISQLConditionBuilder interface {
	Build(db IPersistenceStorage, query interface{}, args ...interface{}) IPersistenceStorage
}

type SQLConditionBuilder struct {
	Type string
}

func (scb *SQLConditionBuilder) Build(db IPersistenceStorage, query interface{}, args ...interface{}) IPersistenceStorage {
	if scb.Type == "or" {
		return db.Or(query, args...)
	}
	return db.Where(query, args...)
}

func NewSQLConditionBuilder(conditionType string) *SQLConditionBuilder {
	return &SQLConditionBuilder{Type: conditionType}
}
