package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/TheJokersThief/frozen-throne/frozen_throne/config"

	gcs "cloud.google.com/go/storage"
)

type GCSStorage struct {
	Context       context.Context
	Config        *config.GCSConfig
	ServiceConfig *config.Config
	Client        *gcs.Client
	Bucket        *gcs.BucketHandle
}

// NewGCSStorage returns a GCS storage engine
func NewGCSStorage(config config.Config, ctx context.Context) *GCSStorage {
	client, err := gcs.NewClient(context.Background())
	if err != nil {
		panic(fmt.Errorf("storage.NewClient: %v", err))
	}

	bkt := client.Bucket(config.GCSConfig.Bucket)

	return &GCSStorage{
		Config:  &config.GCSConfig,
		Client:  client,
		Bucket:  bkt,
		Context: ctx,
	}
}

// PlaceKey writes an object with the name supplied to a GCS bucket
func (s *GCSStorage) PlaceKey(name string) error {
	return s.writeObject(name, name)
}

// RemoveKey deletes an object with the supplied name from a GCS bucket
func (s *GCSStorage) RemoveKey(name string) error {
	objName := fmt.Sprintf("%v/%v", s.Config.BaseFolder, name)
	obj := s.Bucket.Object(objName)
	if err := obj.Delete(s.Context); err != nil {
		return err
	}

	return nil
}

// GetKey retrieves the bucket object and returns its contents
func (s *GCSStorage) GetKey(name string) (string, error) {
	objName := fmt.Sprintf("%v/%v", s.Config.BaseFolder, name)
	obj := s.Bucket.Object(objName)

	r, err := obj.NewReader(s.Context)
	if err != nil {
		return "", err
	}
	defer r.Close()

	slurp, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(slurp), nil
}

// LogAction creates an object with audit data. The intention is to use lifecycle rules to clean up old, unneeded logs.
func (s *GCSStorage) LogAction(action string, user string, repo string, metadata ThroneMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	timestampPrefix := metadata.Timestamp.Format("2006-01-02_15:04:05")
	objName := fmt.Sprintf("%s/%s_%s", s.ServiceConfig.AuditLogKey, timestampPrefix, repo)
	return s.writeObject(objName, string(data))
}

// writeObject is a helper function to write an object to a GCS bucket in the config's base folder
func (s *GCSStorage) writeObject(name string, content string) error {
	objName := fmt.Sprintf("%v/%v", s.Config.BaseFolder, name)
	obj := s.Bucket.Object(objName)

	w := obj.NewWriter(s.Context)
	if _, err := fmt.Fprint(w, name); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}
