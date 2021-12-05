package main

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"frozen_throne/frozen_throne/storage"
)

type Config struct {
	WriteSecret    string `envconfig:"WRITE_SECRET" required:"true"`
	ReadOnlySecret string `envconfig:"READ_ONLY_SECRET" required:"true"`
	StorageMethod  string `envconfig:"STORAGE_METHOD" default:"gcs"`
	AuditLogKey    string `envconfig:"AUDIT_LOG_KEY" default:"audit_log"`
}

type FrozenThrone struct {
	Name    string
	Config  interface{}
	Storage storage.StorageInterface

	storage.GCSConfig
	storage.RedisConfig
}

func NewFrozenThrone(name string) *FrozenThrone {
	config := Config{}

	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	var throneStorage storage.StorageInterface
	switch config.StorageMethod {
	case "gcs":
		throneStorage = storage.NewGCSStorage()
	case "redis":
		throneStorage = storage.NewRedisStorage()
	default:
		panic("Storage method does not match the methods available")
	}

	return &FrozenThrone{
		Name:    name,
		Config:  config,
		Storage: throneStorage,
	}
}

func (f *FrozenThrone) Freeze(name string) (bool, error) {
	return f.Storage.PlaceKey(name)
}

func (f *FrozenThrone) Unfreeze(name string) (bool, error) {
	return f.Storage.RemoveKey(name)
}

func (f *FrozenThrone) Check(name string) (bool, error) {
	return f.Storage.GetKey(name)
}

func IngestHttp(w http.ResponseWriter, r *http.Request) {
	// throne := NewFrozenThrone("test")
}
