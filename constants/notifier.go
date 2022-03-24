package constants

const (
	NotificationIDPrefix      = "nof-"
	NotificationPostTableName = "notification_post"
	PostStatusSending         = "sending"
	PostStatusFailed          = "failed"
	PostStatusFinished        = "finished"
	NotifyTypeEmail           = "email"
	NotifyTypeSMS             = "sms"
	PlatformIAAS              = "iaas"
	PlatformPrivate           = "private_cloud"

	CodeSuccess = 0
)

var NotificationPostColumns = []string{
	"notification_post_id",
	"owner",
	"notify_type",
	"title",
	"content",
	"short_content",
	"status",
	"create_time",
	"status_time",
	"email_address",
}

var ValidStatus = []string{
	PostStatusFailed,
	PostStatusFinished,
	PostStatusSending,
}
