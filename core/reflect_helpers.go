package core

import (
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strconv"
	"strings"
	"time"
)


func SetUpStructField(structF reflect.Value, v interface{}) error {
	switch structF.Kind() {
	case reflect.Int:
		structF.Set(reflect.ValueOf(v))
	case reflect.Int8:
		structF.Set(reflect.ValueOf(v))
	case reflect.Int16:
		structF.Set(reflect.ValueOf(v))
	case reflect.Int32:
		structF.Set(reflect.ValueOf(v))
	case reflect.Int64:
		v := v.(int64)
		if !structF.OverflowInt(v) {
			structF.SetInt(v)
		} else {
			return fmt.Errorf("can't set field with value %d", v)
		}
	case reflect.Uint:
		structF.Set(reflect.ValueOf(v))
	case reflect.Uint8:
		structF.Set(reflect.ValueOf(v))
	case reflect.Uint16:
		structF.Set(reflect.ValueOf(v))
	case reflect.Uint32:
		structF.Set(reflect.ValueOf(v))
	case reflect.Uint64:
		v := v.(uint64)
		if !structF.OverflowUint(v) {
			structF.SetUint(v)
		} else {
			return fmt.Errorf("can't set field with value %d", v)
		}
	case reflect.Bool:
		vI := reflect.ValueOf(v)
		switch vI.Kind() {
		case reflect.String:
			v := v.(string)
			structF.SetBool(v != "")
		case reflect.Bool:
			structF.SetBool(v.(bool))
		}
	case reflect.String:
		v := v.(string)
		structF.SetString(v)
	case reflect.Float32:
		structF.Set(reflect.ValueOf(v))
	case reflect.Float64:
		v := v.(float64)
		structF.SetFloat(v)
	case reflect.Struct:
		switch structF.Interface().(type) {
		case time.Time:
			v := v.(time.Time)
			structF.Set(reflect.ValueOf(v))
		case gorm.DeletedAt:
			v := v.(gorm.DeletedAt)
			structF.Set(reflect.ValueOf(v))
		case ContentType:
			v := v.(ContentType)
			structF.Set(reflect.ValueOf(v))
		default:
			v1 := reflect.ValueOf(v)
			if v1.Kind() == reflect.Ptr {
				structF.Set(v1.Elem())
			} else {
				structF.Set(v1)
			}
		}
	case reflect.Slice:
		switch structF.Type().Name() {
		case "StringArray":
			v := v.(pq.StringArray)
			structF.Set(reflect.ValueOf(v))
		}
	}
	return nil
}

func GetGoMonolithFieldTypeFromGormField(gormField *schema.Field) GoMonolithFieldType {
	var t GoMonolithFieldType
	if gormField.PrimaryKey {
		return PositiveIntegerFieldType
	}
	if gormField.DataType == "json" {
		return JSONFieldType
	}
	if gormField.FieldType.Name() == "StringArray" {
		return ArrayFieldType
	}
	switch gormField.FieldType.Kind() {
	case reflect.Bool:
		t = BooleanFieldType
	case reflect.Int:
		t = IntegerFieldType
	case reflect.Int8:
		t = IntegerFieldType
	case reflect.Int16:
		t = IntegerFieldType
	case reflect.Int32:
		t = IntegerFieldType
	case reflect.Int64:
		t = BigIntegerFieldType
	case reflect.Uint:
		t = PositiveIntegerFieldType
	case reflect.Uint8:
		t = PositiveIntegerFieldType
	case reflect.Uint16:
		t = PositiveIntegerFieldType
	case reflect.Uint32:
		t = PositiveIntegerFieldType
	case reflect.Uint64:
		t = PositiveBigIntegerFieldType
	case reflect.String:
		t = TextFieldType
	case reflect.Float32:
		t = FloatFieldType
	case reflect.Float64:
		t = FloatFieldType
	case reflect.Slice:
		t = ArrayFieldType
	case reflect.Struct:
	}
	return t
}

func IsTruthyValue(value interface{}) bool {
	r := reflect.TypeOf(value)
	if value == nil {
		return false
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		s := reflect.ValueOf(value)
		return s.Len() != 0
	} else if typeString == "bool" {
		return value.(bool) == true
	} else if r.Kind() == reflect.Struct {
	} else if typeString == "string" {
		return value != ""
	} else if typeString == "int" {
		return value.(int) != 0
	} else if typeString == "Month" {
		return value.(int) != 0
	}
	return true
}

func TransformValueForWidget(value interface{}) interface{} {
	r := reflect.TypeOf(value)
	if value == nil {
		return ""
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		if r.Name() == "JSON" {
			return string(value.([]byte))
		}
		newSlice := make([]string, 0)
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			newSlice = append(newSlice, TransformValueForWidget(s.Index(i).Interface()).(string))
		}
		return newSlice
	} else if r.Kind() == reflect.Bool {
		return strconv.FormatBool(value.(bool))
	} else if r.Kind() == reflect.Struct {
		s := reflect.ValueOf(value)
		switch s.Interface().(type) {
		case time.Time:
			return value.(time.Time).Format(CurrentConfig.D.GoMonolith.DateFormat)
		case gorm.DeletedAt:
			return value.(gorm.DeletedAt).Time.Format(CurrentConfig.D.GoMonolith.DateFormat)
		case ContentType:
			ct := value.(ContentType)
			return (&ct).String()
		}
		return ""
	} else if r.Kind() == reflect.Ptr {
		// @todo, handle pointer to time.Time
		s := reflect.Indirect(reflect.ValueOf(value))
		if !s.IsValid() {
			return nil
		}
		switch s.Interface().(type) {
		case time.Time:
			return value.(*time.Time)
		case ContentType:
			ct := value.(ContentType)
			return (&ct).String()
		}
	} else if r.Kind() == reflect.Uint8 {
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	} else if r.Kind() == reflect.Uint16 {
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	} else if r.Kind() == reflect.Uint32 {
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	} else if r.Kind() == reflect.Uint64 {
		return strconv.FormatUint(value.(uint64), 10)
	} else if r.Kind() == reflect.Int8 {
		return strconv.FormatInt(int64(value.(int8)), 10)
	} else if r.Kind() == reflect.Int16 {
		return strconv.FormatInt(int64(value.(int16)), 10)
	} else if r.Kind() == reflect.Int32 {
		return strconv.FormatInt(int64(value.(int32)), 10)
	} else if r.Kind() == reflect.Int64 {
		return strconv.FormatInt(value.(int64), 10)
	} else if r.Kind() == reflect.Float32 {
		return strconv.FormatFloat(float64(value.(float32)), 'E', -1, 32)
	} else if r.Kind() == reflect.Float64 {
		return strconv.FormatFloat(value.(float64), 'E', -1, 64)
	} else if typeString == "string" {
		return value
	} else if typeString == "int" {
		return strconv.Itoa(value.(int))
	} else if typeString == "uint" {
		return fmt.Sprint(value.(uint))
	} else if typeString == "Month" {
		return strconv.Itoa(int(value.(time.Month)))
	}
	return value
}

func TransformDateTimeValueForWidget(value interface{}) interface{} {
	r := reflect.TypeOf(value)
	if value == nil {
		return ""
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		newSlice := make([]string, 0)
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			newSlice = append(newSlice, TransformValueForWidget(s.Index(i).Interface()).(string))
		}
		return newSlice
	} else if r.Kind() == reflect.Bool {
		return strconv.FormatBool(value.(bool))
	} else if r.Kind() == reflect.Struct {
		s := reflect.ValueOf(value)
		switch s.Interface().(type) {
		case time.Time:
			return value.(time.Time).Format(CurrentConfig.D.GoMonolith.DateTimeFormat)
		case gorm.DeletedAt:
			return value.(gorm.DeletedAt).Time.Format(CurrentConfig.D.GoMonolith.DateTimeFormat)
		}
		return ""
	} else if r.Kind() == reflect.Ptr {
		// @todo, handle pointer to time.Time
		s := reflect.Indirect(reflect.ValueOf(value))
		if !s.IsValid() {
			return nil
		}
		switch s.Interface().(type) {
		case time.Time:
			return value.(*time.Time)
		}
	} else if typeString == "string" {
		return value
	} else if typeString == "int" {
		return strconv.Itoa(value.(int))
	} else if typeString == "uint" {
		return fmt.Sprint(value.(uint))
	} else if typeString == "Month" {
		return strconv.Itoa(int(value.(time.Month)))
	}
	return value
}

func TransformValueForOperator(value interface{}) interface{} {
	r := reflect.TypeOf(value)
	if value == nil {
		return ""
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		if r.Name() == "JSON" {
			return string(value.([]byte))
		}
		newSlice := make([]string, 0)
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			newSlice = append(newSlice, TransformValueForWidget(s.Index(i).Interface()).(string))
		}
		return newSlice
	} else if r.Kind() == reflect.Bool {
		return strconv.FormatBool(value.(bool))
	} else if r.Kind() == reflect.Struct {
		s := reflect.ValueOf(value)
		switch s.Interface().(type) {
		case time.Time:
			return value.(time.Time).Format(CurrentConfig.D.GoMonolith.DateFormat)
		case gorm.DeletedAt:
			return value.(gorm.DeletedAt).Time.Format(CurrentConfig.D.GoMonolith.DateFormat)
		}
		return ""
	} else if r.Kind() == reflect.Ptr {
		// @todo, handle pointer to time.Time
		s := reflect.Indirect(reflect.ValueOf(value))
		if !s.IsValid() {
			return nil
		}
		switch s.Interface().(type) {
		case time.Time:
			return value.(*time.Time)
		}
	} else if typeString == "string" {
		boolean, err := strconv.ParseBool(value.(string))
		if err == nil {
			return boolean
		}
		return value
	} else if typeString == "int" {
		return strconv.Itoa(value.(int))
	} else if typeString == "uint" {
		return fmt.Sprint(value.(uint))
	} else if typeString == "Month" {
		return strconv.Itoa(int(value.(time.Month)))
	}
	return value
}

func TransformValueForListDisplay(value interface{}, forExportP ...bool) string {
	forExport := false
	if len(forExportP) > 0 {
		forExport = forExportP[0]
	}
	r := reflect.TypeOf(value)
	if value == nil {
		return ""
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		if r.Name() == "JSON" {
			return string(value.(datatypes.JSON))
		}
		newSlice := make([]string, 0)
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			newSlice = append(newSlice, TransformValueForWidget(s.Index(i).Interface()).(string))
		}
		return strings.Join(newSlice, ",")
	} else if r.Kind() == reflect.Bool {
		if !forExport {
			v := value.(bool)
			if v {
				return "<i class=\"fa fa-check-circle\" aria-hidden=\"TRUE\" style=\"color:green;\"></i>"
			}
			return "<i class=\"fa fa-times-circle\" aria-hidden=\"TRUE\" style=\"color:red;\"></i>"
		}
		return strconv.FormatBool(value.(bool))
	} else if r.Kind() == reflect.Struct {
		s := reflect.ValueOf(value)
		switch s.Interface().(type) {
		case time.Time:
			return value.(time.Time).Format(CurrentConfig.D.GoMonolith.DateFormat)
		case gorm.DeletedAt:
			return value.(gorm.DeletedAt).Time.Format(CurrentConfig.D.GoMonolith.DateFormat)
		}
		return ""
	} else if r.Kind() == reflect.Ptr {
		// @todo, handle pointer to time.Time
		s := reflect.Indirect(reflect.ValueOf(value))
		if !s.IsValid() {
			return ""
		}
		switch s.Interface().(type) {
		case time.Time:
			return value.(*time.Time).String()
		}
	} else if r.Kind() == reflect.Uint8 {
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	} else if r.Kind() == reflect.Uint16 {
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	} else if r.Kind() == reflect.Uint32 {
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	} else if r.Kind() == reflect.Uint64 {
		return strconv.FormatUint(value.(uint64), 10)
	} else if r.Kind() == reflect.Int8 {
		return strconv.FormatInt(int64(value.(int8)), 10)
	} else if r.Kind() == reflect.Int16 {
		return strconv.FormatInt(int64(value.(int16)), 10)
	} else if r.Kind() == reflect.Int32 {
		return strconv.FormatInt(int64(value.(int32)), 10)
	} else if r.Kind() == reflect.Int64 {
		return strconv.FormatInt(value.(int64), 10)
	} else if r.Kind() == reflect.Float32 {
		return strconv.FormatFloat(float64(value.(float32)), 'E', -1, 32)
	} else if r.Kind() == reflect.Float64 {
		return strconv.FormatFloat(value.(float64), 'E', -1, 64)
	} else if typeString == "string" {
		return value.(string)
	} else if typeString == "int64" {
		return strconv.FormatInt(value.(int64), 10)
	} else if typeString == "int" {
		return strconv.Itoa(value.(int))
	} else if typeString == "uint" {
		return fmt.Sprint(value.(uint))
	} else if typeString == "Month" {
		return strconv.Itoa(int(value.(time.Month)))
	}
	return value.(string)
}
