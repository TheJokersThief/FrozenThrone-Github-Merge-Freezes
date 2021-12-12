package config

type Config struct {
	WriteSecret   string `envconfig:"WRITE_SECRET" required:"true"`
	StorageMethod string `envconfig:"STORAGE_METHOD" default:"gcs"`
	AuditLogKey   string `envconfig:"AUDIT_LOG_KEY" default:"audit_log"`

	GCSConfig
	RedisConfig
	FileConfig
}

type GCSConfig struct {
	Bucket     string `envconfig:"GCS_BUCKET"`
	BaseFolder string `envconfig:"GCS_BUCKET_BASEFOLDER" default:""`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"`
	User     string `envconfig:"REDIS_USER"`
	Password string `envconfig:"REDIS_PASSWORD"`
}

type FileConfig struct {
	BaseFolder string `envconfig:"STORAGE_FOLDER" default:"../storage"`
}
