package servers

import (
	config "blcMod/internal/config"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Checking correct of server and it's fields creation
func TestNewServer(t *testing.T) {
	// testing data
	serverConfig := &config.ServerConfiguration{
		Name: "testServer1",
		Adress: url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
			Path:   "/",
		},
		Readiness: true,
	}

	currServ := NewServerFunc(serverConfig.Name, serverConfig)

	if currServ == nil {
		t.Fatal("testNewServer err with server, function NewServer return nil")
	}

	if currServ.Name != serverConfig.Name {
		t.Errorf("testNewServer err with field name, expected: %v, have: %v", serverConfig.Name, currServ.Name)
	}

	if currServ.ServerConfig.Adress.String() != "http://localhost:8080/" {
		t.Errorf("testNewServer err with field adress, expected: localhost:8080/, have: %v", currServ.ServerConfig.Adress.String())
	}

	if currServ.ServerConfig.Readiness != true {
		t.Errorf("testNewServer err with field readiness, expected: true, have: %v", currServ.ServerConfig.Readiness)
	}
}

// Checking if server is correctly started
func TestServerStart(t *testing.T) {
	// testing data
	serverConfig := &config.ServerConfiguration{
		Name: "testServer1",
		Adress: url.URL{
			Host: "localhost:8081",
			Path: "/",
		},
		Readiness: true,
	}
	currServ := NewServerFunc(serverConfig.Name, serverConfig)

	// fast creation of http serv for test
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server %v is working", serverConfig.Name)
	}))
	defer testServer.Close()

	// take test server against real to be more independent
	currServ.server.Handler = testServer.Config.Handler
	currServ.server.Addr = serverConfig.Adress.String()

	// send request to server
	resp, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatalf("TestServerStart err with sending get req, error is %v", err)
	}
	defer resp.Body.Close()

	// check statis code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestServerStart err with status code, expect: %d, have %d", http.StatusOK, resp.StatusCode)
	}

	// reading main body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("TestServerStart err with reading body, error is %v", err)
	}

	// Проверяем содержимое тела ответа
	expectedBody := fmt.Sprintf("Server %s is working", currServ.Name)
	if string(body) != expectedBody {
		t.Errorf("TestServerStart err with body, expect: %s, have %s", expectedBody, string(body))
	}
}
