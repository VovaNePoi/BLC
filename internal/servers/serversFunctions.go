package servers

import (
	config "blcMod/internal/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ServerStruct struct {
	Name         string
	ServerConfig config.ServerConfiguration
	server       *http.Server
}

// Create new server
func NewServerFunc(name string, servConfig *config.ServerConfiguration) *ServerStruct {
	return &ServerStruct{name, *servConfig, &http.Server{Addr: servConfig.Adress.String()}}
}

// starting server
func (s *ServerStruct) Start() {
	var err error

	// Server handler can be in own file and take some data for exmpl.
	// Firstly done simple realization
	http.HandleFunc(fmt.Sprintf("/%s", s.Name), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server %v is working", s.Name)
	})

	fmt.Printf("Starting server on addr: %v", s.server.Addr)
	err = s.server.ListenAndServe()
	if err != nil {
		log.Printf("StartServer err with listening server, error is: %v", err)
	}
}

// Stop server after 10 seconds for doing requests
// wanna add name and time input to close [name] serv after [x minutes]
func (s *ServerStruct) StopServer() error {
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	return s.server.Shutdown(ctx)
}
