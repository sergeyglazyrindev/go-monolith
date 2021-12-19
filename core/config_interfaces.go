package core

import (
	"container/list"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/loads"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

// DBSettings !
type DBSettings struct {
	Type     string `json:"type"` // sqlite, mysql
	Name     string `json:"name"` // File/DB name
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type ConfigOptions struct {
	Theme                  string `yaml:"theme"`
	SiteName               string `yaml:"site_name"`
	ReportingLevel         int    `yaml:"reporting_level"`
	ReportTimeStamp        bool   `yaml:"report_timestamp"`
	DebugDB                bool   `yaml:"debug_db"`
	PageLength             int    `yaml:"page_length"`
	MaxImageHeight         int    `yaml:"max_image_height"`
	MaxImageWidth          int    `yaml:"max_image_width"`
	MaxUploadFileSize      int64  `yaml:"max_upload_file_size"`
	EmailFrom              string `yaml:"email_from"`
	EmailUsername          string `yaml:"email_username"`
	EmailPassword          string `yaml:"email_password"`
	EmailSMTPServer        string `yaml:"email_smtp_server"`
	EmailSMTPServerPort    int    `yaml:"email_smtp_server_port"`
	RootURL                string `yaml:"root_url"`
	RootAdminURL           string `yaml:"root_admin_url"`
	OTPAlgorithm           string `yaml:"otp_algorithm"`
	OTPDigits              int    `yaml:"otp_digits"`
	OTPPeriod              uint   `yaml:"otp_period"`
	OTPSkew                uint   `yaml:"otp_skew"`
	PublicMedia            bool   `yaml:"public_media"`
	RestrictSessionIP      bool   `yaml:"restrict_session_ip"`
	RetainMediaVersions    bool   `yaml:"retain_media_versions"`
	RateLimit              uint   `yaml:"rate_limit"`
	RateLimitBurst         uint   `yaml:"rate_limit_burst"`
	LogHTTPRequests        bool   `yaml:"log_http_requests"`
	HTTPLogFormat          string `yaml:"http_log_format"`
	LogTrail               bool   `yaml:"log_trail"`
	TrailLoggingLevel      int    `yaml:"trail_logging_level"`
	SystemMetrics          bool   `yaml:"system_metrics"`
	UserMetrics            bool   `yaml:"user_metrics"`
	PasswordAttempts       int    `yaml:"password_attempts"`
	PasswordTimeout        int    `yaml:"password_timeout"`
	Logo                   string `yaml:"logo"`
	FavIcon                string `yaml:"fav_icon"`
	AdminCookieName        string `yaml:"admin_cookie_name"`
	APICookieName          string `yaml:"api_cookie_name"`
	SessionDuration        int64  `yaml:"session_duration"`
	SecureCookie           bool   `yaml:"secure_cookie"`
	HTTPOnlyCookie         bool   `yaml:"http_only_cookie"`
	DirectAPISigninByField string `yaml:"direct_api_signin_by_field"`
	PoweredOnSite          string `yaml:"powered_on_site"`
	ForgotCodeExpiration   int    `yaml:"forgot_code_expiration"`
	DateFormat             string `yaml:"date_format"`
	UploadPath             string `yaml:"upload_path"`
	DateTimeFormat         string `yaml:"datetime_format"`
	TimeFormat             string `yaml:"time_format"`
	DateFormatOrder        string `yaml:"date_format_order"`
	AdminPerPage           int    `yaml:"admin_per_page"`
}

type DbOptions struct {
	Default *DBSettings
	Slave   *DBSettings
}

type AuthOptions struct {
	JwtSecretToken    string `yaml:"jwt_secret_token"`
	MinUsernameLength int    `yaml:"min_username_length"`
	MaxUsernameLength int    `yaml:"max_username_length"`
	MinPasswordLength int    `yaml:"min_password_length"`
	SaltLength        int    `yaml:"salt_length"`
}

type AdminOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
	BindIP string `yaml:"bind_ip"`
}

type APIOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
}

type SwaggerOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
	PathToSpec          string `yaml:"path_to_spec"`
	APIEditorListenPort int    `yaml:"api_editor_listen_port"`
}

type ConfigurableConfig struct {
	GoMonolith *ConfigOptions  `yaml:"gomonolith"`
	Test       string          `yaml:"test"`
	Db         *DbOptions      `yaml:"db"`
	Auth       *AuthOptions    `yaml:"auth"`
	Admin      *AdminOptions   `yaml:"admin"`
	API        *APIOptions     `yaml:"api"`
	Swagger    *SwaggerOptions `yaml:"swagger"`
	Debug      bool            `yaml:"debug"`
}

type FieldChoice struct {
	DisplayAs string
	Value     interface{}
}

type IFieldChoiceRegistryInterface interface {
	IsValidChoice(v interface{}) bool
}

type FieldChoiceRegistry struct {
	Choices []*FieldChoice
}

func (fcr *FieldChoiceRegistry) IsValidChoice(v interface{}) bool {
	return false
}

type IFieldFormOptions interface {
	GetName() string
	GetInitial() interface{}
	GetDisplayName() string
	GetValidators() *ValidatorRegistry
	GetChoices() *FieldChoiceRegistry
	GetHelpText() string
	GetWidgetType() string
	GetReadOnly() bool
	GetIsRequired() bool
	GetWidgetPopulate() func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{}
	IsItFk() bool
	GetIsAutocomplete() bool
	GetListFieldWidget() string
}

// Info from config file
type Config struct {
	APISpec                   *loads.Document
	D                         *ConfigurableConfig
	TemplatesFS               embed.FS
	OverridenTemplatesFS      *embed.FS
	LocalizationFS            embed.FS
	CustomLocalizationFS      *embed.FS
	RequiresCsrfCheck         func(c *gin.Context) bool
	PatternsToIgnoreCsrfCheck *list.List
	ErrorHandleFunc           func(int, string, string)
	InTests                   bool
	ConfigContent             []byte
	DebugTests                bool
}

func (c *Config) GetPathToTemplate(templateName string) string {
	return fmt.Sprintf("templates/go-monolith/%s/%s.html", c.D.GoMonolith.Theme, templateName)
}

func (c *Config) GetTemplateContent(templatePath string) []byte {
	templateContent := make([]byte, 0)
	var err error
	if c.OverridenTemplatesFS != nil {
		templateContent, err = c.OverridenTemplatesFS.ReadFile(templatePath)
		if err != nil {
			templateContent, _ = c.TemplatesFS.ReadFile(templatePath)
		}
	} else {
		templateContent, _ = c.TemplatesFS.ReadFile(templatePath)
	}
	return templateContent
}

func (c *Config) GetPathToUploadDirectory() string {
	return fmt.Sprintf("%s/%s", os.Getenv("GOMONOLITH_PATH"), c.D.GoMonolith.UploadPath)
}

func (c *Config) GetURLToUploadDirectory() string {
	return fmt.Sprintf("/%s", c.D.GoMonolith.UploadPath)
}

func (ucc *ConfigurableConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawStuff ConfigurableConfig
	raw := rawStuff{
		Admin: &AdminOptions{BindIP: "0.0.0.0"},
		Auth:  &AuthOptions{SaltLength: 16},
		GoMonolith: &ConfigOptions{
			Theme:             "default",
			SiteName:          "Go-Monolith",
			ReportingLevel:    0,
			ReportTimeStamp:   false,
			DebugDB:           false,
			PageLength:        100,
			MaxImageHeight:    600,
			MaxImageWidth:     800,
			MaxUploadFileSize: int64(25 * 1024 * 1024),
			RootURL:           "/",
			RootAdminURL:      "/admin",
			OTPAlgorithm:      "sha1",
			OTPDigits:         6,
			OTPPeriod:         uint(30),
			OTPSkew:           uint(5),
			PublicMedia:       false,
			//LogDelete: true,
			//LogAdd: true,
			//LogEdit: true,
			//LogRead: false,
			//CacheTranslation: false,
			RestrictSessionIP:   false,
			RetainMediaVersions: true,
			RateLimit:           uint(3),
			RateLimitBurst:      uint(3),
			//APILogRead: false,
			//APILogEdit: true,
			//APILogAdd: true,
			//APILogDelete: true,
			LogHTTPRequests:        true,
			HTTPLogFormat:          "%a %>s %B %U %D",
			LogTrail:               false,
			TrailLoggingLevel:      2,
			SystemMetrics:          false,
			UserMetrics:            false,
			PasswordAttempts:       5,
			PasswordTimeout:        15,
			Logo:                   "/static-inbuilt/go-monolith/logo.png",
			FavIcon:                "/static-inbuilt/go-monolith/favicon.ico",
			AdminCookieName:        "go-monolith-admin",
			APICookieName:          "go-monolith-api",
			SessionDuration:        3600,
			SecureCookie:           false,
			HTTPOnlyCookie:         true,
			DirectAPISigninByField: "username",
			ForgotCodeExpiration:   10,
			DateFormat:             "01/_2/2006",
			DateTimeFormat:         "01/_2/2006 15:04",
			TimeFormat:             "15:04",
			UploadPath:             "uploads",
			DateFormatOrder:        "mm/dd/yyyy",
			AdminPerPage:           10,
		},
	}
	// Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*ucc = ConfigurableConfig(raw)
	return nil

}

var CurrentConfig *Config

// Reads info from config file
func NewConfig(file string) *Config {
	file = os.Getenv("GOMONOLITH_PATH") + "/" + file
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal("Config file is missing: ", file)
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	c := Config{}
	err = yaml.Unmarshal([]byte(content), &c.D)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	c.ConfigContent = content
	c.PatternsToIgnoreCsrfCheck = list.New()
	c.PatternsToIgnoreCsrfCheck.PushBack("/ignorecsrfcheck/")
	c.RequiresCsrfCheck = func(c *gin.Context) bool {
		for e := CurrentConfig.PatternsToIgnoreCsrfCheck.Front(); e != nil; e = e.Next() {
			pathToIgnore := e.Value.(string)
			if c.Request.URL.Path == pathToIgnore {
				return false
			}
		}
		return true
	}
	CurrentConfig = &c
	return &c
}

// Reads info from config file
func NewSwaggerSpec(file string) *loads.Document {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal("Config file is missing: ", file)
	}
	doc, err := loads.Spec(file)
	if err != nil {
		panic(err)
	}
	return doc
}
