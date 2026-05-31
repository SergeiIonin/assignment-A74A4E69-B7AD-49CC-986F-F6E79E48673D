package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/cmd/handlers"
	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/config"
	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	log.Println("Config loaded: %+v", cfg)

	httpClient := NewHTTPClient()

	profileSvc := internal.NewProfileServiceImpl(cfg.ExternalBaseUrl+cfg.UserPath, httpClient)
	todoSvc := internal.NewTodoServiceImpl(cfg.ExternalBaseUrl+cfg.TodosPath+"/user", httpClient)
	dashboardSvc := internal.NewDashboardImpl(profileSvc, todoSvc)

	handlers := handlers.NewHandlers(dashboardSvc)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/dashboard/:id", handlers.ExecutionTimeHandler, handlers.DashboardHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Println("Starting http server")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	time.Sleep(2 * time.Second) // makes sense in scenario of kube-proxy hasn't yet updated iptables and app in the pod is still receiving traffic

	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

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
