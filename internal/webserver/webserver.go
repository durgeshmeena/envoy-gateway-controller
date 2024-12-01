// package webserver

package main

import (
	"net/http"
	"time"

	"github.com/go-logr/logr"
	
	"github.com/durgeshmeena/envoy-gateway-controller/internal/pkg/logging"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/middleware"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/router"
)

// var (
// 	webLog = logging.NewLoggerWithName("WebServer")
// )

func startServer(webLog logr.Logger) error {
	handlers := router.GetRouters(webLog)

	loggedMiddleware := middleware.LoggingMiddleware(webLog)
	loggedHandlers := loggedMiddleware(handlers)

	server := &http.Server{
		Addr:              ":3000",
		Handler:           loggedHandlers,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Listen and Server Web Server
	func() {
		if err := server.ListenAndServe(); err != nil {
			// log.Fatalf("Failed to start server: %v", err)
			webLog.Error(err, "Failed to start server")
		}
	}()

	return nil
}

func main() {
	// logging.InitLogger()
	webLog := logging.InitLogger().WithName("WebServer")
	// log.Println("Starting Web Server")
	webLog.Info("Starting Web Server", "time", time.Now())
	if err := startServer(webLog); err != nil {
		// log.Println("Failed to start server")
		webLog.Error(err, "Failed to start server")
	}
}
