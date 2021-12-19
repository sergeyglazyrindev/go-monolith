package models

/*
	Logging model designed to keep all data about admin actions.
*/

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Action !
type Action int

func (a Action) Read() Action {
	return 1
}

// Added @
func (a Action) Added() Action {
	return 2
}

// Modified !
func (a Action) Modified() Action {
	return 3
}

// Deleted !
func (a Action) Deleted() Action {
	return 4
}

// LoginSuccessful !
func (a Action) LoginSuccessful() Action {
	return 5
}

// LoginDenied !
func (a Action) LoginDenied() Action {
	return 6
}

// Logout !
func (a Action) Logout() Action {
	return 7
}

// PasswordResetRequest !
func (a Action) PasswordResetRequest() Action {
	return 8
}

// PasswordResetDenied !
func (a Action) PasswordResetDenied() Action {
	return 9
}

// PasswordResetSuccessful !
func (a Action) PasswordResetSuccessful() Action {
	return 10
}

// GetSchema !
func (a Action) GetSchema() Action {
	return 11
}

// Custom !
func (a Action) Custom() Action {
	return 99
}

func HumanizeAction(action Action) string {
	switch action {
	case 1:
		return "read"
	case 2:
		return "added"
	case 3:
		return "modified"
	case 4:
		return "deleted"
	case 5:
		return "login successful"
	case 6:
		return "login denied"
	case 7:
		return "logout"
	case 8:
		return "password reset request"
	case 9:
		return "password reset denied"
	case 10:
		return "password reset successful"
	case 11:
		return "read schema"
	default:
		return "unknown"
	}
}

// Log !
type Log struct {
	ID            uint             `gorm:"primarykey"`
	ContentType   core.ContentType `gomonolithform:"ReadonlyField" gomonolith:"list,search"`
	ContentTypeID uint
	ModelPK       uint      `gomonolith:"list,search" gomonolithform:"ReadonlyField"`
	Action        Action    `gomonolithform:"ReadonlyField" gomonolith:"list,search"`
	Username      string    `gomonolithform:"ReadonlyField" gomonolith:"list,search"`
	Activity      string    `gorm:"type:text" gomonolithform:"ReadonlyTextareaFieldOptions" gomonolith:"list,search"`
	CreatedAt     time.Time `gomonolithform:"DateTimeFieldOptions" gomonolith:"list,search"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (l Log) String() string {
	return fmt.Sprintf("Log %s", strconv.Itoa(int(l.ID)))
}

//// Save !
//func (l *Log) Save() {
//	// database.Save(l)
//	//if l.Action == l.Action.Modified() || l.Action == l.Action.Deleted() {
//	//	l.RollBack = preloaded.RootURL + "revertHandler/?log_id=" + fmt.Sprint(l.ID)
//	//}
//	// database.Save(l)
//}

//// ParseRecord !
//func (l *Log) ParseRecord(a reflect.Value, modelName string, ID uint, user *core.User, action Action, r *http.Request) (err error) {
//	//modelName = strings.ToLower(modelName)
//	//model, _ := model2.NewModel(modelName, false)
//	//s, ok := model2.GetSchema(model.Interface())
//	//if !ok {
//	//	errMsg := fmt.Sprintf("Unable to find schema (%s)", modelName)
//	//	debug.Trail(debug.ERROR, errMsg)
//	//	return fmt.Errorf(errMsg)
//	//}
//	//l.Username = user.Username
//	//l.TableName = modelName
//	//l.TableID = int(ID)
//	//l.Action = action
//	//
//	//// Check if the value passed is a pointer
//	//if a.Kind() == reflect.Ptr {
//	//	a = a.Elem()
//	//}
//	//
//	//jsonifyValue := map[string]string{
//	//	"_IP": r.RemoteAddr,
//	//}
//	//for _, f := range s.Fields {
//	//	if !f.IsMethod {
//	//		if f.Type == preloaded.CFK {
//	//			jsonifyValue[f.Name+"ID"] = fmt.Sprint(a.FieldByName(f.Name + "ID").Interface())
//	//		} else if f.Type == preloaded.CDATE {
//	//			val := time.Time{}
//	//			if a.FieldByName(f.Name).Type().Kind() == reflect.Ptr {
//	//				if a.FieldByName(f.Name).IsNil() {
//	//					jsonifyValue[f.Name] = ""
//	//				} else {
//	//					val, _ = a.FieldByName(f.Name).Elem().Interface().(time.Time)
//	//					jsonifyValue[f.Name] = val.Format("2006-01-02 15:04:05 -0700")
//	//				}
//	//
//	//			} else {
//	//				val, _ = a.FieldByName(f.Name).Interface().(time.Time)
//	//				jsonifyValue[f.Name] = val.Format("2006-01-02 15:04:05 -0700")
//	//			}
//	//
//	//		} else {
//	//			jsonifyValue[f.Name] = fmt.Sprint(a.FieldByName(f.Name).Interface())
//	//		}
//	//
//	//	}
//	//}
//	//json1, _ := json.Marshal(jsonifyValue)
//	//l.Activity = string(json1)
//	//
//	return nil
//}

//// SignIn !
//func (l *Log) SignIn(user string, action Action, r *http.Request) (err error) {
//
//	l.Username = user
//	l.Action = action
//	loginStatus := ""
//	if r.Context().Value(preloaded.CKey("login-status")) != nil {
//		loginStatus = r.Context().Value(preloaded.CKey("login-status")).(string)
//	}
//	jsonifyValue := map[string]string{
//		"IP":           r.RemoteAddr,
//		"Login-Status": loginStatus,
//	}
//	for k, v := range r.Header {
//		jsonifyValue[k] = strings.Join(v, ";")
//	}
//
//	json1, _ := json.Marshal(jsonifyValue)
//	l.Activity = string(json1)
//
//	return nil
//}
//
//// PasswordReset !
//func (l *Log) PasswordReset(user string, action Action, r *http.Request) (err error) {
//
//	l.Username = user
//	l.Action = action
//	jsonifyValue := map[string]string{
//		"IP":           r.RemoteAddr,
//		"Reset-Status": r.FormValue("reset-status"),
//	}
//	for k, v := range r.Header {
//		jsonifyValue[k] = strings.Join(v, ";")
//	}
//
//	json1, _ := json.Marshal(jsonifyValue)
//	l.Activity = string(json1)
//
//	return nil
//}
