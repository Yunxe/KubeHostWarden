package probe

import (
	"fmt"
	"kubehostwarden/utils/middleware"
	"log"
	"net/http"
)

func NewServer() {
	var httpServer http.Server

	httpMux := http.NewServeMux()

	httpMux.HandleFunc("/health", health)
	httpMux.HandleFunc("/register",Register)

	httpServer.Handler = middleware.Cors(httpMux)
	httpServer.Addr = ":8080"

	fmt.Printf("Starting Probe Http Server on port %s\n", httpServer.Addr)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start probe http server: %v", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status": "OK"}`))
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
