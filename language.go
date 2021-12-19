package gomonolith

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"os"
)

type LanguageCommand struct {
}

func (c LanguageCommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &core.CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}

	commandRegistry.AddAction("add", &AddLanguageHandler{})
	if len(os.Args) > 2 {
		action = os.Args[2]
		isCorrectActionPassed = commandRegistry.IsRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return nil
	}
	return commandRegistry.RunAction(subaction, "", args)
}

func (c LanguageCommand) GetHelpText() string {
	return "Manipulate languages in the system"
}

type AddLanguageHandlerOptions struct {
	Code string `short:"c" required:"true" description:"Language you wanna add to the system'"`
}

type AddLanguageHandler struct {
}

func (command AddLanguageHandler) Proceed(subaction string, args []string) error {
	var opts = &AddLanguageHandlerOptions{}
	parser := flags.NewParser(opts, flags.Default)
	var err error
	_, err = parser.ParseArgs(args)
	if len(args) == 0 {
		var help string = `
Please provide flag -c which is code of the language
`
		fmt.Printf(help)
		return nil
	}
	if err != nil {
		return err
	}
	langs := [][]string{
		{"Abkhaz", "аҧсуа бызшәа, аҧсшәа", "ab"},
		{"Afar", "Afaraf", "aa"},
		{"Afrikaans", "Afrikaans", "af"},
		{"Akan", "Akan", "ak"},
		{"Albanian", "Shqip", "sq"},
		{"Amharic", "አማርኛ", "am"},
		{"Arabic", "العربية", "ar"},
		{"Aragonese", "aragonés", "an"},
		{"Armenian", "Հայերեն", "hy"},
		{"Assamese", "অসমীয়া", "as"},
		{"Avaric", "авар мацӀ, магӀарул мацӀ", "av"},
		{"Avestan", "avesta", "ae"},
		{"Aymara", "aymar aru", "ay"},
		{"Azerbaijani", "azərbaycan dili", "az"},
		{"Bambara", "bamanankan", "bm"},
		{"Bashkir", "башҡорт теле", "ba"},
		{"Basque", "euskara, euskera", "eu"},
		{"Belarusian", "беларуская мова", "be"},
		{"Bengali, Bangla", "বাংলা", "bn"},
		{"Bihari", "भोजपुरी", "bh"},
		{"Bislama", "Bislama", "bi"},
		{"Bosnian", "bosanski jezik", "bs"},
		{"Breton", "brezhoneg", "br"},
		{"Bulgarian", "български език", "bg"},
		{"Burmese", "ဗမာစာ", "my"},
		{"Catalan", "català", "ca"},
		{"Chamorro", "Chamoru", "ch"},
		{"Chechen", "нохчийн мотт", "ce"},
		{"Chichewa, Chewa, Nyanja", "chiCheŵa, chinyanja", "ny"},
		{"Chinese", "中文 (Zhōngwén), 汉语, 漢語", "zh"},
		{"Chuvash", "чӑваш чӗлхи", "cv"},
		{"Cornish", "Kernewek", "kw"},
		{"Corsican", "corsu, lingua corsa", "co"},
		{"Cree", "ᓀᐦᐃᔭᐍᐏᐣ", "cr"},
		{"Croatian", "hrvatski jezik", "hr"},
		{"Czech", "čeština, český jazyk", "cs"},
		{"Danish", "dansk", "da"},
		{"Divehi, Dhivehi, Maldivian", "ދިވެހި", "dv"},
		{"Dutch", "Nederlands, Vlaams", "nl"},
		{"Dzongkha", "རྫོང་ཁ", "dz"},
		{"Esperanto", "Esperanto", "eo"},
		{"Estonian", "eesti, eesti keel", "et"},
		{"Ewe", "Eʋegbe", "ee"},
		{"Faroese", "føroyskt", "fo"},
		{"Fijian", "vosa Vakaviti", "fj"},
		{"Filipino", "Filipino", "fl"},
		{"Finnish", "suomi, suomen kieli", "fi"},
		{"French", "français, langue française", "fr"},
		{"Fula, Fulah, Pulaar, Pular", "Fulfulde, Pulaar, Pular", "ff"},
		{"Galician", "galego", "gl"},
		{"Georgian", "ქართული", "ka"},
		{"German", "Deutsch", "de"},
		{"Greek (modern)", "ελληνικά", "el"},
		{"Guaraní", "Avañe'ẽ", "gn"},
		{"Gujarati", "ગુજરાતી", "gu"},
		{"Haitian, Haitian Creole", "Kreyòl ayisyen", "ht"},
		{"Hausa", "(Hausa) هَوُسَ", "ha"},
		{"Hebrew (modern)", "עברית", "he"},
		{"Herero", "Otjiherero", "hz"},
		{"Hindi", "हिन्दी, हिंदी", "hi"},
		{"Hiri Motu", "Hiri Motu", "ho"},
		{"Hungarian", "magyar", "hu"},
		{"Interlingua", "Interlingua", "ia"},
		{"Indonesian", "Bahasa Indonesia", "id"},
		{"Interlingue", "Originally called Occidental; then Interlingue after WWII", "ie"},
		{"Irish", "Gaeilge", "ga"},
		{"Igbo", "Asụsụ Igbo", "ig"},
		{"Inupiaq", "Iñupiaq, Iñupiatun", "ik"},
		{"Ido", "Ido", "io"},
		{"Icelandic", "Íslenska", "is"},
		{"Italian", "Italiano", "it"},
		{"Inuktitut", "ᐃᓄᒃᑎᑐᑦ", "iu"},
		{"Japanese", "日本語 (にほんご)", "ja"},
		{"Javanese", "ꦧꦱꦗꦮ, Basa Jawa", "jv"},
		{"Kalaallisut, Greenlandic", "kalaallisut, kalaallit oqaasii", "kl"},
		{"Kannada", "ಕನ್ನಡ", "kn"},
		{"Kanuri", "Kanuri", "kr"},
		{"Kashmiri", "कश्मीरी, كشميري‎", "ks"},
		{"Kazakh", "қазақ тілі", "kk"},
		{"Khmer", "ខ្មែរ, ខេមរភាសា, ភាសាខ្មែរ", "km"},
		{"Kikuyu, Gikuyu", "Gĩkũyũ", "ki"},
		{"Kinyarwanda", "Ikinyarwanda", "rw"},
		{"Kyrgyz", "Кыргызча, Кыргыз тили", "ky"},
		{"Komi", "коми кыв", "kv"},
		{"Kongo", "Kikongo", "kg"},
		{"Korean", "한국어", "ko"},
		{"Kurdish", "Kurdî, كوردی‎", "ku"},
		{"Kwanyama, Kuanyama", "Kuanyama", "kj"},
		{"Latin", "latine, lingua latina", "la"},
		{"Luxembourgish, Letzeburgesch", "Lëtzebuergesch", "lb"},
		{"Ganda", "Luganda", "lg"},
		{"Limburgish, Limburgan, Limburger", "Limburgs", "li"},
		{"Lingala", "Lingála", "ln"},
		{"Lao", "ພາສາລາວ", "lo"},
		{"Lithuanian", "lietuvių kalba", "lt"},
		{"Luba-Katanga", "Tshiluba", "lu"},
		{"Latvian", "latviešu valoda", "lv"},
		{"Manx", "Gaelg, Gailck", "gv"},
		{"Macedonian", "македонски јазик", "mk"},
		{"Malagasy", "fiteny malagasy", "mg"},
		{"Malay", "bahasa Melayu, بهاس ملايو‎", "ms"},
		{"Malayalam", "മലയാളം", "ml"},
		{"Maltese", "Malti", "mt"},
		{"Māori", "te reo Māori", "mi"},
		{"Marathi (Marāṭhī)", "मराठी", "mr"},
		{"Marshallese", "Kajin M̧ajeļ", "mh"},
		{"Mongolian", "Монгол хэл", "mn"},
		{"Nauruan", "Dorerin Naoero", "na"},
		{"Navajo, Navaho", "Diné bizaad", "nv"},
		{"Northern Ndebele", "isiNdebele", "nd"},
		{"Nepali", "नेपाली", "ne"},
		{"Ndonga", "Owambo", "ng"},
		{"Norwegian Bokmål", "Norsk bokmål", "nb"},
		{"Norwegian Nynorsk", "Norsk nynorsk", "nn"},
		{"Norwegian", "Norsk", "no"},
		{"Nuosu", "ꆈꌠ꒿ Nuosuhxop", "ii"},
		{"Southern Ndebele", "isiNdebele", "nr"},
		{"Occitan", "occitan, lenga d'òc", "oc"},
		{"Ojibwe, Ojibwa", "ᐊᓂᔑᓈᐯᒧᐎᓐ", "oj"},
		{"Old Church Slavonic, Church Slavonic, Old Bulgarian", "ѩзыкъ словѣньскъ", "cu"},
		{"Oromo", "Afaan Oromoo", "om"},
		{"Oriya", "ଓଡ଼ିଆ", "or"},
		{"Ossetian, Ossetic", "ирон æвзаг", "os"},
		{"(Eastern) Punjabi", "ਪੰਜਾਬੀ", "pa"},
		{"Pāli", "पाऴि", "pi"},
		{"Persian (Farsi)", "فارسی", "fa"},
		{"Polish", "język polski, polszczyzna", "pl"},
		{"Pashto, Pushto", "پښتو", "ps"},
		{"Portuguese", "Português", "pt"},
		{"Quechua", "Runa Simi, Kichwa", "qu"},
		{"Romansh", "rumantsch grischun", "rm"},
		{"Kirundi", "Ikirundi", "rn"},
		{"Romanian", "Română", "ro"},
		{"Russian", "Русский", "ru"},
		{"Sanskrit (Saṁskṛta)", "संस्कृतम्", "sa"},
		{"Sardinian", "sardu", "sc"},
		{"Sindhi", "सिन्धी, سنڌي، سندھی‎", "sd"},
		{"Northern Sami", "Davvisámegiella", "se"},
		{"Samoan", "gagana fa'a Samoa", "sm"},
		{"Sango", "yângâ tî sängö", "sg"},
		{"Serbian", "српски језик", "sr"},
		{"Scottish Gaelic, Gaelic", "Gàidhlig", "gd"},
		{"Shona", "chiShona", "sn"},
		{"Sinhalese, Sinhala", "සිංහල", "si"},
		{"Slovak", "slovenčina, slovenský jazyk", "sk"},
		{"Slovene", "slovenski jezik, slovenščina", "sl"},
		{"Somali", "Soomaaliga, af Soomaali", "so"},
		{"Southern Sotho", "Sesotho", "st"},
		{"Spanish", "Español", "es"},
		{"Sundanese", "Basa Sunda", "su"},
		{"Swahili", "Kiswahili", "sw"},
		{"Swati", "SiSwati", "ss"},
		{"Swedish", "svenska", "sv"},
		{"Tamil", "தமிழ்", "ta"},
		{"Telugu", "తెలుగు", "te"},
		{"Tajik", "тоҷикӣ, toçikī, تاجیکی‎", "tg"},
		{"Thai", "ไทย", "th"},
		{"Tigrinya", "ትግርኛ", "ti"},
		{"Tibetan Standard, Tibetan, Central", "བོད་ཡིག", "bo"},
		{"Turkmen", "Türkmen, Түркмен", "tk"},
		{"Tagalog", "Wikang Tagalog", "tl"},
		{"Tswana", "Setswana", "tn"},
		{"Tonga (Tonga Islands)", "faka Tonga", "to"},
		{"Turkish", "Türkçe", "tr"},
		{"Tsonga", "Xitsonga", "ts"},
		{"Tatar", "татар теле, tatar tele", "tt"},
		{"Twi", "Twi", "tw"},
		{"Tahitian", "Reo Tahiti", "ty"},
		{"Uyghur", "ئۇيغۇرچە‎, Uyghurche", "ug"},
		{"Ukrainian", "Українська", "uk"},
		{"Urdu", "اردو", "ur"},
		{"Uzbek", "Oʻzbek, Ўзбек, أۇزبېك‎", "uz"},
		{"Venda", "Tshivenḓa", "ve"},
		{"Vietnamese", "Tiếng Việt", "vi"},
		{"Volapük", "Volapük", "vo"},
		{"Walloon", "walon", "wa"},
		{"Welsh", "Cymraeg", "cy"},
		{"Wolof", "Wollof", "wo"},
		{"Western Frisian", "Frysk", "fy"},
		{"Xhosa", "isiXhosa", "xh"},
		{"Yiddish", "ייִדיש", "yi"},
		{"Yoruba", "Yorùbá", "yo"},
		{"Zhuang, Chuang", "Saɯ cueŋƅ, Saw cuengh", "za"},
		{"Zulu", "isiZulu", "zu"},
	}
	database := core.NewDatabaseInstance()
	db := database.Db
	languageExists := db.Model(&core.Language{}).Where("code = ?", opts.Code).Limit(1).Find(&core.Language{})
	if languageExists.RowsAffected > 0 {
		return fmt.Errorf("language %s is already in the system", opts.Code)
	}
	found := false
	for _, lang := range langs {
		if lang[2] != opts.Code {
			continue
		}
		found = true
		l := core.Language{
			EnglishName: lang[0],
			Name:        lang[1],
			Code:        lang[2],
			Active:      false,
		}
		db.Create(&l)
	}
	if !found {
		return fmt.Errorf("language %s doesn't exists", opts.Code)
	}
	core.Trail(core.OK, "Language has been added, active it in admin panel")
	return nil
}

func (command AddLanguageHandler) GetHelpText() string {
	return "Add language"
}
