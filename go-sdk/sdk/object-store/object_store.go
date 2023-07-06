package objectstore

import (
	"fmt"
	regexp2 "regexp"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

type ObjectStore struct {
	logger          logr.Logger
	objStore        nats.ObjectStore
	objectStoreName string
}

func NewObjectStore(logger logr.Logger, jetstream nats.JetStreamContext) (*ObjectStore, error) {
	objectStoreName := viper.GetString("nats.object-store")

	logger = logger.WithName("[OBJECT STORE]")

	objStore, err := initObjectStoreDeps(logger, jetstream, objectStoreName)
	if err != nil {
		return nil, err
	}

	return &ObjectStore{
		logger:          logger,
		objStore:        objStore,
		objectStoreName: objectStoreName,
	}, nil
}

func initObjectStoreDeps(logger logr.Logger, jetstream nats.JetStreamContext,
	objectStoreName string) (nats.ObjectStore, error) {
	if objectStoreName != "" {
		objStore, err := jetstream.ObjectStore(objectStoreName)
		if err != nil {
			return nil, fmt.Errorf("error initializing object store: %w", err)
		}

		return objStore, nil
	}

	logger.Info("Object store not defined. Skipping object store initialization.")
	return nil, nil
}

func (ob ObjectStore) List(regexp ...string) ([]string, error) {
	if ob.objStore == nil {
		return nil, errors.ErrUndefinedObjectStore
	}

	objStoreList, err := ob.objStore.List()
	if err != nil {
		return nil, fmt.Errorf("error listing objects from the object store: %w", err)
	}

	var pattern *regexp2.Regexp

	if len(regexp) > 0 && regexp[0] != "" {
		pat, err := regexp2.Compile(regexp[0])
		if err != nil {
			return nil, fmt.Errorf("error compiling regexp: %w", err)
		}

		pattern = pat
	}

	var response []string

	for _, objName := range objStoreList {
		if pattern == nil || pattern.MatchString(objName.Name) {
			response = append(response, objName.Name)
		}
	}

	//nolint: gomnd
	ob.logger.V(2).Info("Files successfully listed", "Object store", ob.objectStoreName)

	return response, nil
}

func (ob ObjectStore) Get(key string) ([]byte, error) {
	if ob.objStore == nil {
		return nil, errors.ErrUndefinedObjectStore
	}

	response, err := ob.objStore.GetBytes(key)
	if err != nil {
		return nil, fmt.Errorf("error retrieving object with key %s from the object store: %w", key, err)
	}

	//nolint: gomnd
	ob.logger.V(1).Info("File successfully retrieved from object store", key, ob.objectStoreName)

	return response, nil
}

func (ob ObjectStore) Save(key string, payload []byte) error {
	if ob.objStore == nil {
		return errors.ErrUndefinedObjectStore
	}

	if payload == nil {
		return errors.ErrEmptyPayload
	}

	_, err := ob.objStore.PutBytes(key, payload)
	if err != nil {
		return fmt.Errorf("error storing object to the object store: %w", err)
	}

	//nolint: gomnd
	ob.logger.V(1).Info("File successfully stored in object store", key, ob.objectStoreName)

	return nil
}

func (ob ObjectStore) Delete(key string) error {
	if ob.objStore == nil {
		return errors.ErrUndefinedObjectStore
	}

	err := ob.objStore.Delete(key)
	if err != nil {
		return fmt.Errorf("error retrieving object with key %s from the object store: %w", key, err)
	}

	//nolint: gomnd
	ob.logger.V(1).Info("File successfully deleted in object store", key, ob.objectStoreName)

	return nil
}

func (ob ObjectStore) Purge(regexp ...string) error {
	if ob.objStore == nil {
		return errors.ErrUndefinedObjectStore
	}

	var pattern *regexp2.Regexp

	if len(regexp) > 0 && regexp[0] != "" {
		pat, err := regexp2.Compile(regexp[0])
		if err != nil {
			return fmt.Errorf("error compiling regexp: %w", err)
		}

		pattern = pat
	}

	objects, err := ob.List()
	if err != nil {
		return fmt.Errorf("error listing objects from the object store: %w", err)
	}

	for _, objectName := range objects {
		if pattern == nil || pattern.MatchString(objectName) {
			ob.logger.V(1).Info("Deleting object", "Object key", objectName) //nolint: gomnd

			err := ob.objStore.Delete(objectName)
			if err != nil {
				return fmt.Errorf("error purging objects from the object store: %w", err)
			}
		}
	}

	//nolint: gomnd
	ob.logger.V(2).Info("Files successfully purged", "Object Store", ob.objectStoreName)

	return nil
}
