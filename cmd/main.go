package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/cmd/handlers"
	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/config"
	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	fmt.Printf("Config loaded: %+v\n", cfg)

	httpClient := NewHTTPClient()

	profileSvc := internal.NewProfileServiceImpl(cfg.ExternalBaseUrl+cfg.UserPath, httpClient)
	todoSvc := internal.NewTodoServiceImpl(cfg.ExternalBaseUrl+cfg.TodosPath+"/user", httpClient)
	dashboardSvc := internal.NewDashboardImpl(profileSvc, todoSvc)

	handlers := handlers.NewHandlers(dashboardSvc)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/dashboard/:id", handlers.ExecutionTimeHandler, handlers.DashboardHandler)

	fmt.Println("Starting http server")
	r.Run(cfg.Host + ":" + fmt.Sprintf("%d", cfg.Port))
}

func NewHTTPClient() *http.Client {
	// Configure the transport
	transport := &http.Transport{
		// Connection pooling
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,

		// Timeouts for various stages of the connection
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // Connection timeout
			KeepAlive: 30 * time.Second, // Keep-alive period
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
		ResponseHeaderTimeout: 10 * time.Second, // Time to wait for response headers
		ExpectContinueTimeout: 1 * time.Second,  // Time to wait for 100-continue response

		// TLS configuration
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // Enforce TLS 1.2 or higher
		},
	}

	// Create the HTTP client with a global timeout
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Overall request timeout
	}

	return client
}
