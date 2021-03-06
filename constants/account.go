package constants

const (
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
	"lang",
	"email",
	"phone",
	"status",
	"role",
	"currency",
	"gravatar_email",
	"create_time",
	"status_time",
}
