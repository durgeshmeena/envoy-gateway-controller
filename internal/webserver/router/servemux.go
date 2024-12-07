package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/durgeshmeena/envoy-gateway-controller/internal/pkg/utils/datastore"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/handler"
	"github.com/go-logr/logr"
)

type User struct {
	User string `json:"user"`
}

func GetRouters(webLog logr.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	mux.HandleFunc("/user/{username}", func(w http.ResponseWriter, r *http.Request) {
		userName := r.PathValue("username")
		fmt.Fprintf(w, "Hello %s", userName)
	})

	mux.HandleFunc("POST /user", func(w http.ResponseWriter, r *http.Request) {
		var user User
		webLog.Info("Received POST request", "time", time.Now())
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			webLog.Error(err, "Failed to decode request body")
			return
		}
		// fmt.Fprintf(w, "User: %s", user.User)
		webLog.Info("User", "user", user.User)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user.User))
	})

	// handle create btp
	fileStore := datastore.NewFileStore(webLog.WithName("FileStore"))
	btpHandler := handler.NewCreateClientBTPHandler(&webLog, fileStore)
	mux.Handle("POST /btp/create", btpHandler)

	return mux
}
