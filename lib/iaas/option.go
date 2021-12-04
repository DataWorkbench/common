package iaas

type Operation struct {
	accessKeyId     string
	secretAccessKey string
}

type Option func(op *Operation)

// WithAccessKey used to send request by specified ak/sk
func WithAccessKey(accessKeyId string, secretAccessKey string) Option {
	return func(op *Operation) {
		op.accessKeyId = accessKeyId
		op.secretAccessKey = secretAccessKey
	}
}
