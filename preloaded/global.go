package preloaded

//// Constants
//
//// cPOST post
//const CPOST = "POST"
//
//// cID true
//const CID = "id"
//
//// cTRUE true
//const CTRUE = "true"
//
//// cJPG jpg
//const CJPG = "jpg"
//
//// cJPEG jpeg
//const CJPEG = "jpeg"
//
//// cPNG png
//const CPNG = "png"
//
//// cGIF gif
//const CGIF = "gif"
//
//// cSTRING !
//const CSTRING = "string"
//
//// cNUMBER !
//const CNUMBER = "number"
//
//// cDATE !
//const CDATE = "date"
//
//// cBOOL !
//const CBOOL = "bool"
//
//// cLIST !
//const CLIST = "list"
//
//// cIMAGE !
//const CIMAGE = "image"
//
//// cFK !
//const CFK = "fk"
//
//// cLINK !
//const CLINK = "link"
//
//// cMONRY !
//const CMONEY = "money"
//
//// cCODE !
//const CCODE = "code"
//
//// cHTML !
//const CHTML = "html"
//
//// cMULTILINGUAL !
//const CMULTILINGUAL = "multilingual"
//
//// cPROGRESSBAR !
//const CPROGRESSBAR = "progress_bar"
//
//// cPASSWORD !
//const CPASSWORD = "password"
//
//// cFILE !
//const CFILE = "file"
//
//// cEMAIL !
//const CEMAIL = "email"
//
//// cM2M !
//const CM2M = "m2m"
//
//// Version number as per Semantic Versioning 2.0.0 (semver.org)
//const Version = "0.6.2"
//
//// VersionCodeName is the cool name we give to versions with significant changes.
//// This name should always be a bug's name starting from A-Z them revolving back.
//// This started at version 0.5.0 (Atlas Moth)
//// 0.6.0 Beetle
//const VersionCodeName = "Beetle"
//
//// Public Global Variables
//
//// Theme is the name of the theme used in GoMonolith.
//var Theme = "default"
//
//// SiteName is the name of the website that shows on title and dashboard.
//var SiteName = "GoMonolith"
//
//// DebugDB prints all SQL statements going to DB.
//var DebugDB = false
//
//// PageLength is the list view max number of records.
//var PageLength = 100
//
//// MaxImageHeight sets the maximum height of an image.
//var MaxImageHeight = 600
//
//// MaxImageWidth sets the maximum width of an image.
//var MaxImageWidth = 800
//
//// MaxUploadFileSize is the maximum upload file size in bytes.
//var MaxUploadFileSize = int64(25 * 1024 * 1024)
//
//// BindIP is the IP the application listens to.
//var BindIP = ""
//
//// EmailFrom identifies where the email is coming from.
//var EmailFrom string
//
//// EmailUsername sets the username of an email.
//var EmailUsername string
//
//// EmailPassword sets the password of an email.
//var EmailPassword string
//
//// EmailSMTPServer sets the name of the SMTP Server in an email.
//var EmailSMTPServer string
//
//// EmailSMTPServerPort sets the port number of an SMTP Server in an email.
//var EmailSMTPServerPort int
//
//// RootURL is where the listener is mapped to.
//var RootURL = "/"
//
//// OTPAlgorithm is the hashing algorithm of OTP.
//var OTPAlgorithm = "sha1"
//
//// OTPDigits is the number of digits for the OTP.
//var OTPDigits = 6
//
//// OTPPeriod is the number of seconds for the OTP to change.
//var OTPPeriod = uint(30)
//
//// OTPSkew is the number of minutes to search around the OTP.
//var OTPSkew = uint(5)
//
//// PublicMedia allows public access to media handler without authentication.
//var PublicMedia = false
//
//// LogDelete adds a log when a record is deleted.
//var LogDelete = true
//
//// LogAdd adds a log when a record is added.
//var LogAdd = true
//
//// LogEdit adds a log when a record is edited.
//var LogEdit = true
//
//// LogRead adds a log when a record is read.
//var LogRead = false
//
//// CacheTranslation allows a translation to store data in a cache memory.
//var CacheTranslation = false
//
//// DefaultMediaPermission is the default permission applied to to files uploaded to the system
//var DefaultMediaPermission = os.FileMode(0644)
//
//// ErrorHandleFunc is a function that will be called everytime Trail is called. It receives
//// one parameter for error level, one for error message and one for runtime stack trace
//var ErrorHandleFunc func(int, string, string)
//
//// AllowedIPs is a list of allowed IPs to access GoMonolith interfrace in one of the following formats:
//// - "*" = Allow all
//// - "" = Allow none
//// - "192.168.1.1" Only allow this IP
//// - "192.168.1.0/24" Allow all IPs from 192.168.1.1 to 192.168.1.254
//// You can also create a list of the above formats using comma to separate them.
//// For example: "192.168.1.1,192.168.1.2,192.168.0.0/24"
//var AllowedIPs = "*"
//
//// BlockedIPs is a list of blocked IPs from accessing GoMonolith interfrace in one of the following formats:
//// - "*" = Block all
//// - "" = Block none
//// - "192.168.1.1" Only block this IP
//// - "192.168.1.0/24" Block all IPs from 192.168.1.1 to 192.168.1.254
//// You can also create a list of the above formats using comma to separate them.
//// For example: "192.168.1.1,192.168.1.2,192.168.0.0/24"
//var BlockedIPs = ""
//
//// RestrictSessionIP is to block access of a user if their IP changes from their original IP during login
//var RestrictSessionIP = false
//
//// RetainMediaVersions is to allow the system to keep files uploaded even after they are changed.
//// This allows the system to "Roll Back" to an older version of the file.
//var RetainMediaVersions = true

//// RateLimit is the maximum number of requests/second for any unique IP
//var RateLimit int64 = 3
//
//// RateLimitBurst is the maximum number of requests for an idle user
//var RateLimitBurst int64 = 3

//// OptimizeSQLQuery selects columns during rendering a form a list to visible fields.
//// This means during the filtering of a form the select statement will not include
//// any field with `hidden` tag. For list it will not select any field with `list_exclude`
//var OptimizeSQLQuery = false
//
//// APILogRead controls the data API's logging for read commands.
//var APILogRead = false
//
//// APILogEdit controls the data API's logging for edit commands.
//var APILogEdit = true
//
//// APILogAdd controls the data API's logging for add commands.
//var APILogAdd = true
//
//// APILogDelete controls the data API's logging for delete commands.
//var APILogDelete = true
//
//// APILogSchema controls the data API's logging for schema commands.
//var APILogSchema = true
//
//// APIPublicRead controls the data API’s public for add commands.
//var APIPublicRead = false
//
//// APIPublicEdit controls the data API's public for edit commands.
//var APIPublicEdit = false
//
//// APIPublicAdd controls the data API's public for add commands.
//var APIPublicAdd = false
//
//// APIPublicDelete controls the data API's public for delete commands.
//var APIPublicDelete = false
//
//// APIPublicSchema controls the data API's public for schema commands.
//var APIPublicSchema = false
//
//// APIDisabledRead controls the data API’s disabled for add commands.
//var APIDisabledRead = false
//
//// APIDisabledEdit controls the data API's disabled for edit commands.
//var APIDisabledEdit = false
//
//// APIDisabledAdd controls the data API's disabled for add commands.
//var APIDisabledAdd = false
//
//// APIDisabledDelete controls the data API's disabled for delete commands.
//var APIDisabledDelete = false
//
//// APIDisabledSchema controls the data API's disabled for schema commands.
//var APIDisabledSchema = false
//
//// APIPreQueryRead controls the data API’s pre query for add commands.
//var APIPreQueryRead = false
//
//// APIPreQueryEdit controls the data API's pre query for edit commands.
//var APIPreQueryEdit = false
//
//// APIPreQueryAdd controls the data API's pre query for add commands.
//var APIPreQueryAdd = false
//
//// APIPreQueryDelete controls the data API's pre query for delete commands.
//var APIPreQueryDelete = false
//
//// APIPreQuerySchema controls the data API's pre query for schema commands.
//var APIPreQuerySchema = false
//
//// APIPostQueryRead controls the data API’s post query for add commands.
//var APIPostQueryRead = false
//
//// APIPostQueryEdit controls the data API's post query for edit commands.
//var APIPostQueryEdit = false
//
//// APIPostQueryAdd controls the data API's post query for add commands.
//var APIPostQueryAdd = false
//
//// APIPostQueryDelete controls the data API's post query for delete commands.
//var APIPostQueryDelete = false
//
//// APIPostQuerySchema controls the data API's post query for schema commands.
//var APIPostQuerySchema = false
//
//// LogHTTPRequests logs http requests to syslog
//var LogHTTPRequests = true
//
///*
//HTTPLogFormat is the format used to log HTTP access
//%a: Client IP address
//%{remote}p: Client port
//%A: Server hostname/IP
//%{local}p: Server port
//%U: Path
//%c: All coockies
//%{NAME}c: Cookie named 'NAME'
//%{GET}f: GET request parameters
//%{POST}f: POST request parameters
//%B: Response length
//%>s: Response code
//%D: Time taken in microseconds
//%T: Time taken in seconds
//%I: Request length
//*/
//var HTTPLogFormat = "%a %>s %B %U %D"
//
//// LogTrail stores Trail logs to syslog
//var LogTrail = false

//// CacheSessions allows GoMonolith to store sessions data in memory
//var CacheSessions = true

//// CachePermissions allows GoMonolith to store permissions data in memory
//var CachePermissions = true
//
//// PasswordAttempts is the maximum number of invalid password attempts before
//// the IP address is blocked for some time from usig the system
//var PasswordAttempts = 5
//
//// PasswordTimeout is the amount of time in minutes the IP will be blocked for after
//// reaching the the maximum invalid password attempts
//var PasswordTimeout = 15
//
//// AllowedHosts is a comma seprated list of allowed hosts for the server to work. The
//// default value if only for development and production domain should be added before
//// deployment
//var AllowedHosts = "0.0.0.0,127.0.0.1,localhost,::1"
//
//// Logo is the main logo that shows on GoMonolith UI
//var Logo = "/static/go-monolith/logo.png"
//
//// FavIcon is the fav icon that shows on GoMonolith UI
//var FavIcon = "/static/go-monolith/favicon.ico"
//
//// Private Global Variables
//// Regex
//var MatchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
//var MatchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
//
//// langMapCache is a computer memory used for storage of frequently or recently used translations.
//var LangMapCache = map[string][]byte{}
//
//var Registered = false
//
//var HandlersRegistered = false
//
//var DefaultProgressBarColor = "#07c"
//
//var SettingsSynched = false

//type ListData struct {
//	Rows  [][]interface{}
//	Count int
//}
//
//// CKey is the standard key used in GoMonolith for context keys
//type CKey string
