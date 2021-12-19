---
sidebar_position: 1
---

# Config

Config is a general config for gomonolith. It has following fields:
```go
// Info from config file
type Config struct {
	// would be used to handle swagger, but maybe it would be refactored later.
	APISpec                   *loads.Document
	// this is data loaded from the config
	D                         *ConfigurableConfig
	// this is embedded template fs
	TemplatesFS               embed.FS
	// this is embedded localization fs
	LocalizationFS            embed.FS
	// used to ignore CSRF checks for admin functinality
	RequiresCsrfCheck         func(c *gin.Context) bool
	// patterns to ignore csrf check
	PatternsToIgnoreCsrfCheck *list.List
	// error handle func that could be used to send important error for example to sentry
	ErrorHandleFunc           func(int, string, string)
	// this is a special variable that is set it to true during tests
	InTests                   bool
	// you can use ConfigContent in your blueprint to initialize blueprint specific configs.
	ConfigContent			  []byte
	// debug tests mode for database, simplifies test debugging
	DebugTests                bool
}
```
An example of the config in the configs/sqlite.yml
```yml
db:
  default:
    type: "sqlite"
    name: "/persist/test.db"
admin:
  listen_port: 8080
api:
  listen_port: 5000
auth:
  max_username_length: 40
  min_username_length: 8
  min_password_length: 8
swagger:
  path_to_spec: configs/api-spec.yml
  listen_port: 8082
  api_editor_listen_port: 8083
gomonolith:
  secure_cookie: false
  http_only_cookie: false
  debug_tests: false
  upload_path: upload-for-tests
```
And there's a lot of options that could be configured through config file.  
You may provide default database alias and slave. This is typical configuration to have a cluster with slaves and one master. So it should suit most of the real use cases.  
You can customize listen_port for admin, etc. Please check out core/config_interfaces.go.  
```go
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
	// theme, should be a subfolder in the gomonolith folder in the TemplatesFS
	Theme                  string `yaml:"theme"`
	// site name that will be used in the admin panel.
	SiteName               string `yaml:"site_name"`
	// currently not widely used, maybe removed later.
	ReportingLevel         int    `yaml:"reporting_level"`
	// currently not widely used, maybe removed later.
	ReportTimeStamp        bool   `yaml:"report_timestamp"`
	// currently not widely used, maybe removed later.
	DebugDB                bool   `yaml:"debug_db"`
	// currently not widely used, maybe removed later.
	PageLength             int    `yaml:"page_length"`
	// to be used during cropping, currently not in use cause cropping is a plannned feature
	MaxImageHeight         int    `yaml:"max_image_height"`
	// to be used during cropping, currently not in use cause cropping is a plannned feature
	MaxImageWidth          int    `yaml:"max_image_width"`
	// max upload file size
	MaxUploadFileSize      int64  `yaml:"max_upload_file_size"`
	EmailFrom              string `yaml:"email_from"`
	EmailUsername          string `yaml:"email_username"`
	EmailPassword          string `yaml:"email_password"`
	EmailSMTPServer        string `yaml:"email_smtp_server"`
	EmailSMTPServerPort    int    `yaml:"email_smtp_server_port"`
	// root url for your API
	RootURL                string `yaml:"root_url"`
	// root url for your admin panel
	RootAdminURL           string `yaml:"root_admin_url"`
	// currently not widely used, maybe removed later.
	OTPAlgorithm           string `yaml:"otp_algorithm"`
	// currently not widely used, maybe removed later.
	OTPDigits              int    `yaml:"otp_digits"`
	// currently not widely used, maybe removed later.
	OTPPeriod              uint   `yaml:"otp_period"`
	// currently not widely used, maybe removed later.
	OTPSkew                uint   `yaml:"otp_skew"`
	// currently not widely used, maybe removed later.
	PublicMedia            bool   `yaml:"public_media"`
	// currently not widely used, maybe removed later.
	RestrictSessionIP      bool   `yaml:"restrict_session_ip"`
	// currently not widely used, maybe removed later.
	RetainMediaVersions    bool   `yaml:"retain_media_versions"`
	// currently not widely used, maybe removed later.
	RateLimit              uint   `yaml:"rate_limit"`
	// currently not widely used, maybe removed later.
	RateLimitBurst         uint   `yaml:"rate_limit_burst"`
	// currently not widely used, maybe removed later.
	LogHTTPRequests        bool   `yaml:"log_http_requests"`
	// currently not widely used, maybe removed later.
	HTTPLogFormat          string `yaml:"http_log_format"`
	// currently not widely used, maybe removed later.
	LogTrail               bool   `yaml:"log_trail"`
	// currently not widely used, maybe removed later.
	TrailLoggingLevel      int    `yaml:"trail_logging_level"`
	// currently not widely used, maybe removed later.
	SystemMetrics          bool   `yaml:"system_metrics"`
	// currently not widely used, maybe removed later.
	UserMetrics            bool   `yaml:"user_metrics"`
	// currently not widely used, maybe removed later.
	PasswordTimeout        int    `yaml:"password_timeout"`
	// currently not widely used, maybe removed later.
	PasswordAttempts       int    `yaml:"password_attempts"`
	// logo for admin panel
	Logo                   string `yaml:"logo"`
	// favicon for admin panel
	FavIcon                string `yaml:"fav_icon"`
	// admin cookie name for user session
	AdminCookieName        string `yaml:"admin_cookie_name"`
	// currently not widely used, maybe removed later.
	APICookieName          string `yaml:"api_cookie_name"`
	SessionDuration        int64  `yaml:"session_duration"`
	SecureCookie           bool   `yaml:"secure_cookie"`
	HTTPOnlyCookie         bool   `yaml:"http_only_cookie"`
	// to be used during signin, it could be email or username.
	DirectAPISigninByField string `yaml:"direct_api_signin_by_field"`
	// currently not widely used, maybe removed later.
	DebugTests             bool   `yaml:"debug_tests"`
	// currently not widely used, maybe removed later.
	PoweredOnSite          string `yaml:"powered_on_site"`
	// used to generate code for password reset
	ForgotCodeExpiration   int    `yaml:"forgot_code_expiration"`
	DateFormat             string `yaml:"date_format"`
	UploadPath             string `yaml:"upload_path"`
	DateTimeFormat         string `yaml:"datetime_format"`
	TimeFormat             string `yaml:"time_format"`
	// order of the fields for your project for date format. used in one form widget.
	DateFormatOrder        string `yaml:"date_format_order"`
	// number of records on one admin list page.
	AdminPerPage           int    `yaml:"admin_per_page"`
}

type DbOptions struct {
	Default *DBSettings
	Slave *DBSettings
}

type AuthOptions struct {
	// currently not widely used, maybe removed later.
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
```
Default values are following:
```go
func (ucc *ConfigurableConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawStuff ConfigurableConfig
	raw := rawStuff{
		Admin: &AdminOptions{BindIP: "0.0.0.0"},
		Auth:  &AuthOptions{SaltLength: 16},
		GoMonolith: &ConfigOptions{
			Theme:             "default",
			SiteName:          "Go Monolith",
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
			Logo:                   "/static-inbuilt/gomonolith/logo.png",
			FavIcon:                "/static-inbuilt/gomonolith/favicon.ico",
			AdminCookieName:        "go-monolith-admin",
			APICookieName:          "go-monolith-api",
			SessionDuration:        3600,
			SecureCookie:           false,
			HTTPOnlyCookie:         true,
			DirectAPISigninByField: "username",
			DebugTests:             false,
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
```
