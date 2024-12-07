package webserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/durgeshmeena/envoy-gateway-controller/internal/pkg/logging"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/middleware"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/router"
)

// implement Runnable interface, so that it can be aaded into the controller
type WebServer struct{}

// Start implements the Runnable interface
func (w *WebServer) Start(ctx context.Context) error {
	return StartServer(ctx)
}

// verify that WebServer implements the Runnable interface
var _ manager.Runnable = &WebServer{}

// StartServer starts the web server
func StartServer(ctx context.Context) error {
	webLog := logging.Log.WithName("WebServer")

	handlers := router.GetRouters(webLog)
	loggedMiddleware := middleware.LoggingMiddleware(webLog)
	loggedHandlers := loggedMiddleware(handlers)

	server := &http.Server{
		Addr:              ":3000",
		Handler:           loggedHandlers,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		webLog.Info("Starting Web Server", "port", 3000)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			webLog.Error(err, "Failed to start server")
			errChan <- err
		}
        // If the error is http.ErrServerClosed, it's expected; do not log as error
	}()

	// Wait for the context to be done
	select {
	case <-ctx.Done():
		webLog.Info("Shutting down Web Server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			webLog.Error(err, "Web server forced to shutdown")
			return err
		}
		webLog.Info("Web server stopped gracefully")
		return nil

	case err := <-errChan:
		webLog.Error(err, "Web server encountered an error")
		return err
	}
}