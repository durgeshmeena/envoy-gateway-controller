package datastore

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	model "github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/model"
)

// store data in file

// func readJSONFile(filepath string) ([]byte, error) {
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	byteValue, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return byteValue, nil
// }

// func convertJSONToYAML(jsonData []byte) ([]byte, error) {
// 	// convert json to yaml
// 	yamlData, err := yaml.JSONToYAML(jsonData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return yamlData, nil
// }

// func main() {
// 	// get path of the file
// 	cmd, err := os.Getwd()
// 	if err != nil {
// 		log.Println("Error getting current working directory: ", err)
// 	}
// 	fileName := "input.json"
// 	fileDir := "internal/pkg/utils/datastore"
// 	filePath := filepath.Join(cmd, fileDir, fileName)
// 	fmt.Println("Input File path: ", filePath)

// 	jsonDataBytes, err := os.ReadFile(filePath)
// 	if err != nil {
// 		log.Println("Error reading JSON file: ", err)
// 	}

// 	var btp model.BTP
// 	if err := json.Unmarshal(jsonDataBytes, &btp); err != nil {
// 		log.Println("Error unmarshalling JSON data: ", err)
// 	}

// 	if err := saveBTPToFile(btp); err != nil {
// 		log.Println("Error saving BTP to file: ", err)
// 	}
// 	log.Println("BTP saved to file successfully")
// }

// read existing BTPs from file
func readBTPsFromFile(filePath string) ([]model.BTP, error) {
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
func SaveBTPToFile(btp model.BTP) error {
	// get path of the file
	cmd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory: ", err)
	}
	fileName := "btps.json"
	fileDir := "internal/pkg/utils/datastore"
	filePath := filepath.Join(cmd, fileDir, fileName)
	fmt.Println("Store File path: ", filePath)

	// check if file exists, if not create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist, creating file")
		file, err := os.Create(filePath)
		if err != nil {
			log.Println("Error creating file: ", err)
			return err
		}
		file.Close()
	}

	//  read existing BTPs from file
	btps, err := readBTPsFromFile(filePath)
	if err != nil {
		log.Println("Error reading BTPs from file: ", err)
		return err
	}

	// append new BTP to existing BTPs
	btps = append(btps, btp)

	// convert BTPs to JSON,
	// marshal entire slice of BTPs to JSON
	btpsJSON, err := json.MarshalIndent(btps, "", "    ")
	if err != nil {
		log.Println("Error marshalling BTPs to JSON: ", err)
		return err
	}

	// write JSON to file
	err = os.WriteFile(filePath, btpsJSON, 0644)
	if err != nil {
		log.Println("Error writing BTPs to file: ", err)
		return err
	}

	return nil
}
