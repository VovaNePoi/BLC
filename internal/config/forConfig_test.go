package config

import (
	"net/url"
	"os"
	"reflect"
	"testing"
)

// for testing functionality
func TestWriteToConfig(t *testing.T) {
	// testing data
	testServers := ServersConfigList{
		Servers: map[string]ServerConfiguration{
			"server1": {
				Name:      "Server One",
				Adress:    url.URL{Scheme: "http", Host: "localhost:8080"},
				Readiness: true,
			},
			"server2": {
				Name:      "Server Two",
				Adress:    url.URL{Scheme: "https", Host: "example.com"},
				Readiness: false,
			},
		},
	}

	// fileName for testing
	tempFileName := "test_config.json"

	// Checking simple functions
	err := testServers.WriteToConfig(tempFileName)
	if err != nil {
		t.Fatalf("testingWriteToConfig error with writting data, err is %v", err)
	}

	_, err = os.Stat(tempFileName)
	if os.IsNotExist(err) {
		t.Errorf("testingWriteToConfig error with file creation, err is %v", err)
	}
	defer os.Remove(tempFileName)

	// testing empty input
	emptyServers := ServersConfigList{Servers: map[string]ServerConfiguration{}}
	err = emptyServers.WriteToConfig(tempFileName)
	if err == nil {
		t.Errorf("testingWriteToConfig error with emtyServer must return error, now return nil")
	}

	// testing empty filename
	err = testServers.WriteToConfig()
	if err != nil {
		t.Fatalf("testingWriteToConfig error with writting data without filename, err is %v", err)
	}
	defer os.Remove("config.json") // DANGER!!! Can delete file with all server's info about it.
}

// testing ReadConfig()
func TestGetConfig(t *testing.T) {
	// testing data
	testServers := ServersConfigList{
		Servers: map[string]ServerConfiguration{
			"server1": {
				Name:      "Server One",
				Adress:    url.URL{Scheme: "http", Host: "localhost:8080"},
				Readiness: true,
			},
		},
	}

	// fileName for testing
	tempFileName := "test_read_config.json"

	err := testServers.WriteToConfig(tempFileName)
	if err != nil {
		t.Fatalf("testingReadConfig error with writting data, err is %v", err)
	}
	defer os.Remove(tempFileName)

	// simple tests
	readServers, err := (&ServersConfigList{}).GetConfig(tempFileName)
	if err != nil {
		t.Fatalf("testingReadConfig error with reading data, err is %v", err)
	}

	// checkig if data is correctly read
	if !reflect.DeepEqual(readServers.Servers, testServers.Servers) {
		t.Errorf("testingReadConfig error with reading data, expected: %+v, have: %+v", testServers.Servers, readServers.Servers)
	}

	// testing no xeisting file
	_, err = (&ServersConfigList{}).GetConfig("nonexistent_file.json")
	if err == nil {
		t.Errorf("testingReadConfig error with reading file, err is %v, but must return error", err)
	}

	// try read file "config.json" without input fileName
	// need to create file firstly
	testServers.WriteToConfig()
	readServers, err = (&ServersConfigList{}).GetConfig() // Используем ReadConfig без имени файла
	if err != nil {
		t.Fatalf("testingReadConfig error with reading data from default filename, err is %v", err)
	}
	if !reflect.DeepEqual(readServers.Servers, testServers.Servers) {
		t.Errorf("testingReadConfig error with reading data from default filename, expected: %+v, have: %+v", testServers.Servers, readServers.Servers)
	}
	defer os.Remove("config.json") // DANGER!!! Can delete file with all server's info about it.
}
