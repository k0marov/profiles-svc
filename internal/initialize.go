package internal

import (
	"fmt"
	ory "github.com/ory/client-go"
	"log"
	"net/http"
	"os"
)

func InitializeAndStart() {
	oryProxyPort := os.Getenv("ORY_PROXY_PORT")
	if oryProxyPort == "" {
		oryProxyPort = "4000"
	}

	// register a new Ory client with the URL set to the Ory CLI Proxy
	// we can also read the URL from the env or a config file
	oryCfg := ory.NewConfiguration()
	oryCfg.Servers = ory.ServerConfigurations{{URL: fmt.Sprintf("http://localhost:%s/.ory", oryProxyPort)}}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	repo := NewMongoProfileRepository()
	svc := NewProfileService(repo)
	srv := NewServer(svc, ory.NewAPIClient(oryCfg))
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
