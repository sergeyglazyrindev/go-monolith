package models

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
)

// SettingCategory is a category for system settings
type SettingCategory struct {
	core.Model
	Name string `gomonolith:"list"`
	Icon string `gomonolithform:"ImageFormOptions" `
}

func (m *SettingCategory) String() string {
	return fmt.Sprintf("Setting category %s", m.Name)
}

// DataType is a list of data types used for settings
type DataType int

// String is a type
func (DataType) String() DataType {
	return 1
}

// Integer is a type
func (DataType) Integer() DataType {
	return 2
}

// Float is a type
func (DataType) Float() DataType {
	return 3
}

// Boolean is a type
func (DataType) Boolean() DataType {
	return 4
}

// File is a type
func (DataType) File() DataType {
	return 5
}

// Image is a type
func (DataType) Image() DataType {
	return 6
}

// DateTime is a type
func (DataType) DateTime() DataType {
	return 7
}

func HumanizeDataType(dataType DataType) string {
	switch dataType {
	case 1:
		return "string"
	case 2:
		return "integer"
	case 3:
		return "float"
	case 4:
		return "boolean"
	case 5:
		return "file"
	case 6:
		return "image"
	case 7:
		return "datetime"
	default:
		return "unknown"
	}
}

func DataTypeFromString(dataType string) DataType {
	switch dataType {
	case "string":
		return DataType(1)
	case "integer":
		return DataType(2)
	case "float":
		return DataType(3)
	case "boolean":
		return DataType(4)
	case "file":
		return DataType(5)
	case "image":
		return DataType(6)
	case "datetime":
		return DataType(7)
	default:
		return DataType(0)
	}
}

// Setting model stored system settings
type Setting struct {
	core.Model
	Name         string          `gomonolith:"list,search" gomonolithform:"RequiredFieldOptions"`
	Value        string          `gomonolith:"list" gomonolithform:"DynamicTypeFieldOptions"`
	DefaultValue string          `gomonolith:"list" gomonolithform:"DynamicTypeFieldOptions"`
	DataType     DataType        `gomonolith:"list,search" gomonolithform:"RequiredSelectFieldOptions"`
	Help         string          `gomonolith:"list,search" sql:"type:text;"`
	Category     SettingCategory `gomonolith:"list,search" gomonolithform:"FkRequiredFieldOptions"`
	CategoryID   uint
	Code         string `gomonolith:"search" gomonolithform:"ReadonlyField"`
}

func (s *Setting) GetRealWidget() core.IWidget {
	switch s.DataType {
	case 1:
		return &core.TextWidget{}
	case 2:
		return &core.NumberWidget{}
	case 3:
		widget := &core.NumberWidget{}
		widget.SetAttr("step", "0.1")
		return widget
	case 4:
		return &core.CheckboxWidget{}
	case 5:
		return &core.FileWidget{}
	case 6:
		widget := &core.FileWidget{}
		widget.SetAttr("accept", "image/*")
		return widget
	case 7:
		return &core.DateTimeWidget{}
	default:
		return &core.Widget{}
	}
}

func (s *Setting) String() string {
	return fmt.Sprintf("Setting %s", s.Name)
}

//// Save overides save
//func (s *Setting) Save() {
//	// @todo, probably use it
//	//database.Preload(s)
//	//s.Code = strings.Replace(s.Category.Name, " ", "", -1) + "." + strings.Replace(s.Name, " ", "", -1)
//	//s.ApplyValue()
//	//database.Save(s)
//}

//// ParseFormValue takes the value of a setting from an HTTP request and saves in the instance of setting
//func (s *Setting) ParseFormValue(v []string) {
//	switch s.DataType {
//	case s.DataType.Boolean():
//		tempV := len(v) == 1 && v[0] == "on"
//		if tempV {
//			s.Value = "1"
//		} else {
//			s.Value = "0"
//		}
//	case s.DataType.DateTime():
//		if len(v) == 1 && v[0] != "" {
//			s.Value = v[0] + ":00"
//		} else {
//			s.Value = ""
//		}
//	default:
//		if len(v) == 1 && v[0] != "" {
//			s.Value = v[0]
//		} else {
//			s.Value = ""
//		}
//	}
//}
//
//// GetValue returns an interface representing the value of the setting
//func (s *Setting) GetValue() interface{} {
//	var err error
//	var v interface{}
//
//	switch s.DataType {
//	case s.DataType.String():
//		if s.Value == "" {
//			v = s.DefaultValue
//		} else {
//			v = s.Value
//		}
//	case s.DataType.Integer():
//		if s.Value != "" {
//			v, err = strconv.ParseInt(s.Value, 10, 64)
//			v = int(v.(int64))
//		}
//		if err != nil {
//			v, err = strconv.ParseInt(s.DefaultValue, 10, 64)
//		}
//		if err != nil {
//			v = 0
//		}
//	case s.DataType.Float():
//		if s.Value != "" {
//			v, err = strconv.ParseFloat(s.Value, 64)
//		}
//		if err != nil {
//			v, err = strconv.ParseFloat(s.DefaultValue, 64)
//		}
//		if err != nil {
//			v = 0.0
//		}
//	case s.DataType.Boolean():
//		if s.Value != "" {
//			v = s.Value == "1"
//		}
//		if v == nil {
//			v = s.DefaultValue == "1"
//		}
//	case s.DataType.File():
//		if s.Value == "" {
//			v = s.DefaultValue
//		} else {
//			v = s.Value
//		}
//	case s.DataType.Image():
//		if s.Value == "" {
//			v = s.DefaultValue
//		} else {
//			v = s.Value
//		}
//	case s.DataType.DateTime():
//		if s.Value != "" {
//			v, err = time.Parse("2006-01-02 15:04:05", s.Value)
//		}
//		if err != nil {
//			v, err = time.Parse("2006-01-02 15:04:05", s.DefaultValue)
//		}
//		if err != nil {
//			v = time.Now()
//		}
//	}
//	return v
//}

// @todo, analyze later
//// ApplyValue changes GoMonolith global variables' value based in the setting value
//func (s *Setting) ApplyValue() {
//	v := s.GetValue()
//
//	switch s.Code {
//	case "uAdmin.Theme":
//		preloaded.Theme = strings.Replace(v.(string), "/", "_", -1)
//		preloaded.Theme = strings.Replace(preloaded.Theme, "\\", "_", -1)
//		preloaded.Theme = strings.Replace(preloaded.Theme, "..", "_", -1)
//	case "uAdmin.SiteName":
//		preloaded.SiteName = v.(string)
//	case "uAdmin.ReportingLevel":
//		utils.ReportingLevel = v.(int)
//	case "uAdmin.ReportTimeStamp":
//		utils.ReportTimeStamp = v.(bool)
//	case "uAdmin.DebugDB":
//		if preloaded.DebugDB != v.(bool) {
//			preloaded.DebugDB = v.(bool)
//		}
//	case "uAdmin.PageLength":
//		preloaded.PageLength = v.(int)
//	case "uAdmin.MaxImageHeight":
//		preloaded.MaxImageHeight = v.(int)
//	case "uAdmin.MaxImageWidth":
//		preloaded.MaxImageWidth = v.(int)
//	case "uAdmin.MaxUploadFileSize":
//		preloaded.MaxUploadFileSize = int64(v.(int))
//	case "uAdmin.Port":
//		// Port = v.(int)
//	case "uAdmin.EmailFrom":
//		preloaded.EmailFrom = v.(string)
//	case "uAdmin.EmailUsername":
//		preloaded.EmailUsername = v.(string)
//	case "uAdmin.EmailPassword":
//		preloaded.EmailPassword = v.(string)
//	case "uAdmin.EmailSMTPServer":
//		preloaded.EmailSMTPServer = v.(string)
//	case "uAdmin.EmailSMTPServerPort":
//		preloaded.EmailSMTPServerPort = v.(int)
//	case "uAdmin.RootURL":
//		preloaded.RootURL = v.(string)
//	case "uAdmin.OTPAlgorithm":
//		preloaded.OTPAlgorithm = v.(string)
//	case "uAdmin.OTPDigits":
//		preloaded.OTPDigits = v.(int)
//	case "uAdmin.OTPPeriod":
//		preloaded.OTPPeriod = uint(v.(int))
//	case "uAdmin.OTPSkew":
//		preloaded.OTPSkew = uint(v.(int))
//	case "uAdmin.PublicMedia":
//		preloaded.PublicMedia = v.(bool)
//	case "uAdmin.LogDelete":
//		preloaded.LogDelete = v.(bool)
//	case "uAdmin.LogAdd":
//		preloaded.LogAdd = v.(bool)
//	case "uAdmin.LogEdit":
//		preloaded.LogEdit = v.(bool)
//	case "uAdmin.LogRead":
//		preloaded.LogRead = v.(bool)
//	case "uAdmin.CacheTranslation":
//		preloaded.CacheTranslation = v.(bool)
//	case "uAdmin.AllowedIPs":
//		preloaded.AllowedIPs = v.(string)
//	case "uAdmin.BlockedIPs":
//		preloaded.BlockedIPs = v.(string)
//	case "uAdmin.RestrictSessionIP":
//		preloaded.RestrictSessionIP = v.(bool)
//	case "uAdmin.RetainMediaVersions":
//		preloaded.RetainMediaVersions = v.(bool)
//	case "uAdmin.RateLimit":
//		if preloaded.RateLimit != int64(v.(int)) {
//			preloaded.RateLimit = int64(v.(int))
//			utils.RateLimitMap = map[string]int64{}
//		}
//	case "uAdmin.RateLimitBurst":
//		preloaded.RateLimitBurst = int64(v.(int))
//	case "uAdmin.OptimizeSQLQuery":
//		preloaded.OptimizeSQLQuery = v.(bool)
//	case "uAdmin.APILogRead":
//		preloaded.APILogRead = v.(bool)
//	case "uAdmin.APILogEdit":
//		preloaded.APILogEdit = v.(bool)
//	case "uAdmin.APILogAdd":
//		preloaded.APILogAdd = v.(bool)
//	case "uAdmin.APILogDelete":
//		preloaded.APILogDelete = v.(bool)
//	case "uAdmin.APILogSchema":
//		preloaded.APILogSchema = v.(bool)
//	case "uAdmin.LogHTTPRequests":
//		preloaded.LogHTTPRequests = v.(bool)
//	case "uAdmin.HTTPLogFormat":
//		preloaded.HTTPLogFormat = v.(string)
//	case "uAdmin.LogTrail":
//		preloaded.LogTrail = v.(bool)
//	case "uAdmin.TrailLoggingLevel":
//		interfaces.TrailLoggingLevel = v.(int)
//	case "uAdmin.SystemMetrics":
//		metrics.SystemMetrics = v.(bool)
//	case "uAdmin.UserMetrics":
//		metrics.UserMetrics = v.(bool)
//	case "uAdmin.CacheSessions":
//		preloaded.CacheSessions = v.(bool)
//		if preloaded.CacheSessions {
//			// @todo, probably fix
//			// sessionmodel.LoadSessions()
//		}
//	case "uAdmin.CachePermissions":
//		preloaded.CachePermissions = v.(bool)
//		if preloaded.CachePermissions {
//			// @todo, probably fix
//			// usermodel.LoadPermissions()
//		}
//	case "uAdmin.PasswordAttempts":
//		preloaded.PasswordAttempts = v.(int)
//	case "uAdmin.PasswordTimeout":
//		preloaded.PasswordTimeout = v.(int)
//	case "uAdmin.AllowedHosts":
//		preloaded.AllowedHosts = v.(string)
//	case "uAdmin.Logo":
//		preloaded.Logo = v.(string)
//	case "uAdmin.FavIcon":
//		preloaded.FavIcon = v.(string)
//	}
//}
