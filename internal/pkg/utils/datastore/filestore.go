package datastore

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	model "github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/model"
	"github.com/go-logr/logr"
)

// store data in file

// FileStore struct
type FileStore struct {
	dataStoreLogger logr.Logger
}

// NewFileStore creates a new FileStore instance with logger
func NewFileStore(logger logr.Logger) *FileStore {
	return &FileStore{
		dataStoreLogger: logger,
	}
}

func (fs *FileStore) GetDataStoreFilePath() (string, error) {
	// get path of the file

	fs.dataStoreLogger.Info("Getting dataStore json file path")

	cmd, err := os.Getwd()
	if err != nil {
		fs.dataStoreLogger.Error(err, "Error getting current working directory")
	}
	fileName := "btps.json"
	fileDir := "internal/pkg/utils/datastore"
	filePath := filepath.Join(cmd, fileDir, fileName)
	fs.dataStoreLogger.Info("Store File path: ", "path", filePath)

	// check if file exists, if not create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fs.dataStoreLogger.Info("File does not exist, creating file")
		file, err := os.Create(filePath)
		if err != nil {
			fs.dataStoreLogger.Error(err, "Error creating file")
			return "", err
		}
		file.Close()
	}

	return filePath, nil

}

// read existing BTPs from file
func (fs *FileStore) ReadBTPsFromFile(filePath string) ([]model.BTP, error) {
	var btps []model.BTP

	// read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file: ", err)
		return nil, err
	}

	// if file is empty, return empty slice
	if len(data) == 0 {
		return btps, nil
	}

	// unmarshal JSON data to slice of BTPs
	if err := json.Unmarshal(data, &btps); err != nil {
		log.Println("Error unmarshalling JSON data: ", err)
		return nil, err
	}

	return btps, nil
}

// save BTP to file with validation
func (fs *FileStore) SaveBTPToFile(btp model.BTP) error {

	// get path of the file
	filePath, err := fs.GetDataStoreFilePath()
	if err != nil {
		fs.dataStoreLogger.Error(err, "Error getting dataStore json file path")
		return err
	}

	fs.dataStoreLogger.Info("Reading data from file")
	//  read existing BTPs from file
	btps, err := fs.ReadBTPsFromFile(filePath)
	if err != nil {
		fs.dataStoreLogger.Error(err, "Error reading BTPs from file")
		return err
	}

	// append new BTP to existing BTPs
	btps = append(btps, btp)

	// convert BTPs to JSON,
	// marshal entire slice of BTPs to JSON
	btpsJSON, err := json.MarshalIndent(btps, "", "    ")
	if err != nil {
		fs.dataStoreLogger.Error(err, "Error marshalling BTPs to JSON")
		return err
	}

	// write JSON to file
	err = os.WriteFile(filePath, btpsJSON, 0644)
	if err != nil {
		fs.dataStoreLogger.Error(err, "Error writing BTPs to file")
		return err
	}

	fs.dataStoreLogger.Info("data updated successfully with new BTP")
	return nil
}
