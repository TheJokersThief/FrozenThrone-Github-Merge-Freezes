package frozen_throne

import (
	"context"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/TheJokersThief/frozen-throne/frozen_throne/config"
	"github.com/TheJokersThief/frozen-throne/frozen_throne/storage"
)

type FrozenThrone struct {
	Config  config.Config
	Storage storage.StorageInterface
}

func NewFrozenThrone(ctx context.Context) *FrozenThrone {
	config := config.Config{}

	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	var throneStorage storage.StorageInterface
	switch config.StorageMethod {
	case "gcs":
		throneStorage = storage.NewGCSStorage(config, ctx)
	default:
		panic("Storage method chosen does not match the methods available")
	}

	return &FrozenThrone{
		Config:  config,
		Storage: throneStorage,
	}
}

// Freeze will prevent merges for a repo
func (f *FrozenThrone) Freeze(name string, user string) error {
	err := f.Storage.PlaceKey(name)

	action := "freeze"
	metadata := storage.ThroneMetadata{
		Action:    action,
		User:      user,
		Repo:      name,
		Timestamp: time.Now(),
	}
	f.Storage.LogAction(action, user, name, metadata)
	return err
}

// Thaw will allow merges for a repo
func (f *FrozenThrone) Thaw(name string, user string) error {
	err := f.Storage.RemoveKey(name)

	action := "thaw"
	metadata := storage.ThroneMetadata{
		Action:    action,
		User:      user,
		Repo:      name,
		Timestamp: time.Now(),
	}
	f.Storage.LogAction(action, user, name, metadata)
	return err
}

// Check will see if a repo is currently frozen
func (f *FrozenThrone) Check(name string) (string, error) {
	return f.Storage.GetKey(name)
}
