package core

type Language struct {
	Model
	Code           string `gorm:"uniqueIndex;not null" gomonolithform:"ReadonlyField" gomonolith:"list,search"`
	Name           string `gomonolithform:"ReadonlyField" gomonolith:"list,search"`
	EnglishName    string `gomonolithform:"ReadonlyField" gomonolith:"list,search"`
	Active         bool   `gorm:"default:false" gomonolith:"list"`
	Flag           string `gomonolithform:"ImageFormOptions"`
	RTL            bool   `gorm:"default:false"`
	Default        bool   `gorm:"default:false"`
	AvailableInGui bool   `gorm:"default:false" gomonolithform:"ReadonlyField" gomonolith:"list"`
}

// String !
func (l Language) String() string {
	return l.Code
}

//// Save !
//func (l *Language) Save() {
//	//if l.Default {
//	//	database.Update([]Language{}, "default", false, "\"default\" = ?", true)
//	//	defaultLang = l
//	//}
//	//database.Save(l)
//	//tempActiveLangs := []Language{}
//	//dialect1 := dialect.GetDialectForDb("default")
//	//dialect1.Equals("active", true)
//	//database.Filter(&tempActiveLangs, dialect1.ToString(), true)
//	//ActiveLangs = tempActiveLangs
//	//
//	//tanslationList := []translation.Translation{}
//	//for i := range ActiveLangs {
//	//	tanslationList = append(tanslationList, translation.Translation{
//	//		Active:  ActiveLangs[i].Active,
//	//		Default: ActiveLangs[i].Default,
//	//		Code:    ActiveLangs[i].Code,
//	//		Name:    fmt.Sprintf("%s (%s)", ActiveLangs[i].Name, ActiveLangs[i].EnglishName),
//	//	})
//	//}
//	//
//	//for modelName := range model.Schema {
//	//	for i := range model.Schema[modelName].Fields {
//	//		if model.Schema[modelName].Fields[i].Type == preloaded.CMULTILINGUAL {
//	//			// @todo, redo
//	//			// model.Schema[modelName].Fields[i].Translations = tanslationList
//	//		}
//	//	}
//	//}
//}
