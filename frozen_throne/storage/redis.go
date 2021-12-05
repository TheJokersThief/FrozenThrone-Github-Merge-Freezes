package storage

type RedisStorage struct {
	StorageInterface
}

func NewRedisStorage() *RedisStorage {
	return &RedisStorage{}
}
