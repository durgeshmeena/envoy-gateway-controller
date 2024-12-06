package webserver

import (
	"net/http"
	"time"

	"github.com/durgeshmeena/envoy-gateway-controller/internal/pkg/logging"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/middleware"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/router"
)

// StartServer starts the web server
func StartServer() error {
	webLog := logging.Log.WithName("WebServer")

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

	// print log about starting server
	webLog.Info("Starting Web Server", "port", 3000, "time", time.Now())

	// Listen and Server Web Server
	func() {
		if err := server.ListenAndServe(); err != nil {
			// log.Fatalf("Failed to start server: %v", err)
			webLog.Error(err, "Failed to start server")
		}
	}()

	return nil
}

// func main() {
// 	// logging.InitLogger()
// 	webLog := logging.InitLogger().WithName("WebServer")
// 	// log.Println("Starting Web Server")
// 	webLog.Info("Starting Web Server", "time", time.Now())
// 	if err := StartServer(webLog); err != nil {
// 		// log.Println("Failed to start server")
// 		webLog.Error(err, "Failed to start server")
// 	}
// }
