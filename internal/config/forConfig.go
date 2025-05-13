package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// In this type file saving data about server
type ServerConfiguration struct {
	Name      string  `json:"name"`
	Adress    url.URL `json:"adress"`
	Readiness bool    `json:"readiness"`
}

// Only buffer, don't go to file. In file saved ServerConfiguration
type ServersConfigList struct {
	Servers map[string]ServerConfiguration
}

// write all info about each server to file JSON. If empty input, write to
// "config.json". If file isn't created will create file
func (sl *ServersConfigList) WriteToConfig(fName ...string) error {
	var err error
	var fileName string
	// checking if fName empty
	if len(fName) == 0 {
		fileName = ValidFileName("")
	} else {
		// taking only first name of fName[]
		fileName = ValidFileName(fName[0])
	}

	// check if list empty
	if len(sl.Servers) == 0 {
		return fmt.Errorf("writeToConfig err with data, error is: no data in serverList")
	}

	marshaledData, err := json.Marshal(sl.Servers)
	if err != nil {
		return fmt.Errorf("writeToConfig err with encoding, error is: %v", err)
	}

	// created for closing file
	var flForClose *os.File
	//c check if file created, if not, create file
	if _, err = os.Stat(fileName); os.IsNotExist(err) {
		flForClose, err = os.Create(fileName)
		if err != nil {
			return fmt.Errorf("writeToConfig err with creating file, error is: %v", err)
		}
	}
	defer flForClose.Close()

	err = os.WriteFile(fileName, marshaledData, 0644)
	if err != nil {
		return fmt.Errorf("writeToConfig err with writing to file, error is: %v", err)
	}

	return nil
}

// read all from file(if empty filename read from "config.json") and return ServersList adress
func (sl *ServersConfigList) GetConfig(fName ...string) (*ServersConfigList, error) {
	var err error
	var fileName string
	// checking if fName empty
	if len(fName) == 0 {
		fileName = ValidFileName("")
	} else {
		// checking fName[] and take only first
		fileName = ValidFileName(fName[0])
	}

	slBytes, err := os.ReadFile(fileName)
	if err != nil {
		return &ServersConfigList{}, fmt.Errorf("ReadConfig err with read data from file, error is %v", err)
	}

	err = json.Unmarshal(slBytes, &sl.Servers)
	if err != nil {
		return &ServersConfigList{}, fmt.Errorf("ReadConfig err with deserialize data, error is %v", err)
	}

	return sl, err
}

// checking if name is valid. If empty return "config.json", if no suffix ".json" add it.
// If this tests ok, return fName
func ValidFileName(fName string) string {
	if fName == "" {
		return "config.json"
	} else if !strings.HasSuffix(fName, ".json") {
		return fName + ".json"
	} else {
		return fName
	}
}

// In future can add delete from file, update file and etc.
