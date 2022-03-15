package constants

const (
	Account                        = "account"
	QingcloudSource                = "qingcloud"
	LocalSource                    = "local"
	AccessKeyTableName             = "access_key"
	UserTableName                  = "user"
	QingcloudAccessKeyStatusActive = "active"
)

var AccessKeyColumns = []string{
	"access_key_id",
	"access_key_name",
	"secret_access_key",
	"description",
	"owner",
	"status",
	"create_time",
	"status_time",
	"ip_white_list",
}

var UserColumns = []string{
	"user_id",
	"user_name",
	"password",
	"lang",
	"email",
	"phone",
	"status",
	"role",
	"currency",
	"gravatar_email",
	"privilege",
	"zones",
	"regions",
	"create_time",
	"status_time",
}

// redis
const (
	UserPrefix          = UserTableName
	RedisSeparator      = ":"
	QingcloudUserPrefix = Account + RedisSeparator + QingcloudSource + RedisSeparator + UserPrefix + RedisSeparator
	LocalUserPrefix     = Account + RedisSeparator + LocalSource + RedisSeparator + UserPrefix + RedisSeparator
	DefaultUserPrefix   = LocalUserPrefix
	SessionPrefix       = "session"

	UserCacheBaseSeconds         = 300
	UserCacheRandomSeconds       = 120
	NotExistResourceCacheSeconds = 30

	AccessKeyCacheBaseSeconds = 30
	SessionCacheSeconds       = 3600 * 24
)

// user
const (
	IdPrefixUser         = "usr-"
	IdInstanceUser int64 = iota + 1
	// user status
	UserStatusActive = "active"
	UserStatusBanned = "banned"
	UserStatusDelete = "deleted"
	// access key
	AccessKeyIdLength      = 20
	SecretKeyLength        = 50
	AccessKeyIdLetters     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SecretKeyLetters       = "0123456789abcdefghijklmnopqrstuvwxyz"
	AccessKeyStatusEnable  = "enable"
	AccessKeyStatusDisable = "disable"
	SessionLength          = 50
	SessionLetters         = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ProviderEnfi = "enfi"
)
