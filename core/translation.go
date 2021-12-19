package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
	"sync"
)

// GetLanguage returns the language of the request
func GetLanguage(c *gin.Context) *Language {
	langCookie, err := c.Cookie("language")
	if err != nil || langCookie == "" {
		return GetDefaultLanguage()
	}
	var lang Language
	database := NewDatabaseInstance()
	db := database.Db
	db.Model(Language{}).Where(&Language{Code: langCookie}).First(&lang)
	database.Close()
	return &lang
}

// GetDefaultLanguage returns the default language
func GetDefaultLanguage() *Language {
	if defaultLang != nil {
		return defaultLang
	}
	var lang Language
	database := NewDatabaseInstance()
	defer database.Close()
	db := database.Db
	db.Model(Language{}).Where(&Language{Default: true}).First(&lang)
	defaultLang = &lang
	return &lang
}

// GetActiveLanguages returns a list of active langages
func GetActiveLanguages() []Language {
	var langs []Language
	database := NewDatabaseInstance()
	defer database.Close()
	db := database.Db
	db.Model(Language{}).Where(&Language{Active: true}).Find(&langs)
	return langs
}

// DefaultLang is the default language of the system.
var defaultLang *Language

type translationLoaded map[string]string

var langMapCache map[string]translationLoaded

func readLocalization(languageCode string) translationLoaded {
	langMap, ok := langMapCache[languageCode]
	if ok {
		return langMap
	}
	ret := make(translationLoaded)
	langFile, err := CurrentConfig.LocalizationFS.ReadFile(fmt.Sprintf("localization/%s.json", languageCode))
	if err != nil {
		if CurrentConfig.CustomLocalizationFS != nil {
			langFile, err := CurrentConfig.CustomLocalizationFS.ReadFile(fmt.Sprintf("localization/%s.json", languageCode))
			if err == nil {
				err = json.Unmarshal(langFile, &ret)
			}
		}
	} else {
		err = json.Unmarshal(langFile, &ret)
		if err != nil {
			Trail(ERROR, "Unable to unmarshal json file with language (%s)", err)
		} else {
			if CurrentConfig.CustomLocalizationFS != nil {
				langFile, err := CurrentConfig.CustomLocalizationFS.ReadFile(fmt.Sprintf("localization/%s.json", languageCode))
				if err == nil {
					ret1 := make(translationLoaded)
					err = json.Unmarshal(langFile, &ret1)
					if err == nil {
						for term, translated := range ret1 {
							ret[term] = translated
						}
					}
				}
			}
		}
	}
	langMapCache[languageCode] = ret
	return ret
}

const translateMe = "Translate me ---> "

// SafeCounter is safe to use concurrently.
type NotTranslatedDataImplementation struct {
	mu sync.Mutex
	D  map[string][]string
}

func (ntdi *NotTranslatedDataImplementation) Store(lang string, term string) {
	ntdi.mu.Lock()
	if ntdi.D[lang] == nil {
		ntdi.D[lang] = make([]string, 0)
	}
	ntdi.D[lang] = append(ntdi.D[lang], term)
	ntdi.mu.Unlock()
}

var NotTranslatedData *NotTranslatedDataImplementation

func init() {
	NotTranslatedData = &NotTranslatedDataImplementation{
		D: make(map[string][]string),
	}
}

// @todo, redo
// Tf is a function for translating strings into any given language
// Parameters:
// ===========
//   path (string): This is where to get the translation from. It is in the
//                  format of "GROUPNAME/FILENAME" for example: "gomonolith/system"
//   lang (string): Is the language code. If empty string is passed we will use
//                  the default language.
//   term (string): The term to translate.
//   args (...interface{}): Is a list of args to fill the term with place holders
func Tf(lang string, iTerm interface{}, args ...interface{}) string {
	term := ""
	iTermReflectV := reflect.ValueOf(iTerm)
	httpErrorResponse := false
	itemV := &HTTPErrorResponse{}
	if iTermReflectV.Kind() == reflect.String {
		term = iTerm.(string)
	} else if itemV1, ok := iTerm.(*HTTPErrorResponse); ok {
		httpErrorResponse = true
		itemV = itemV1
		term = itemV1.Code
	} else if itermE, ok := iTerm.(error); ok {
		term = itermE.Error()
	}
	if lang == "" {
		lang = GetDefaultLanguage().Code
	}

	// Check if the path if for an existing model schema
	if langMapCache == nil {
		langMapCache = make(map[string]translationLoaded)
	}
	langMap, ok := langMapCache[lang]
	if !ok {
		langMap = readLocalization(lang)
	}
	// If the term exists, then return it
	if val, ok := langMap[term]; ok {
		if httpErrorResponse && len(itemV.Params) > 0 {
			return fmt.Sprintf(strings.TrimPrefix(val, translateMe), itemV.Params...)
		}
		return strings.TrimPrefix(val, translateMe)
	}
	// If it doesn't exist then add it to the file
	if lang != "en" {
		if CurrentConfig.D.Debug {
			NotTranslatedData.Store(lang, term)
		}
		if httpErrorResponse {
			langMap[term] = translateMe + itemV.Message
			if len(itemV.Params) > 0 {
				return fmt.Sprintf(itemV.Message, itemV.Params...)
			}
		} else {
			langMap[term] = translateMe + term
		}
		Trail(WARNING, "Unknown term %s", term)
		return translateMe + term
	}
	if httpErrorResponse {
		langMap[term] = itemV.Message
	} else {
		langMap[term] = term
	}
	if httpErrorResponse && len(itemV.Params) > 0 {
		return fmt.Sprintf(itemV.Message, itemV.Params...)
	}
	return term
}

func Translate(c *gin.Context, raw string, lang string, args ...bool) string {
	var langParser map[string]json.RawMessage
	err := json.Unmarshal([]byte(raw), &langParser)
	if err != nil {
		return raw
	}
	transtedStr := string(langParser[lang])

	if len(transtedStr) > 2 {
		return transtedStr[1 : len(transtedStr)-1]
	}
	if len(args) > 0 && !args[0] {
		return ""
	}
	language := GetLanguage(c)
	transtedStr = string(langParser[language.Code])

	if len(transtedStr) > 2 {
		return transtedStr[1 : len(transtedStr)-1]
	}
	return ""
}
