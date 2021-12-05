package storage

import (
	gcs "cloud.google.com/go/storage"
)

type GCSStorage struct {
	StorageInterface
	Client gcs.Client
}

type GCSConfig struct {
	Bucket     string `envconfig:"GCS_BUCKET"`
	BaseFolder string `envconfig:"GCS_BUCKET_BASEFOLDER" default:""`
}

func NewGCSStorage() *GCSStorage {
	return &GCSStorage{}
}

func (s *GCSStorage) PlaceKey(name string) (bool, error) {

}

func (s *GCSStorage) RemoveKey(name string) (bool, error) {

}

func (s *GCSStorage) GetKey(name string) (bool, error) {

}

func (s *GCSStorage) LogAction(action string, user string, repo string, metadata map[string]interface{}) (bool, error) {

}
