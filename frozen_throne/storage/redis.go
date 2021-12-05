package storage

type RedisStorage struct {
	StorageInterface
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"`
	User     string `envconfig:"REDIS_USER"`
	Password string `envconfig:"REDIS_PASSWORD"`
}

func NewRedisStorage() *RedisStorage {
	return &RedisStorage{}
}
