package opscenter

import (
	"fmt"
	"kubehostwarden/opscenter/alarm"
	"kubehostwarden/opscenter/probe"
	"kubehostwarden/opscenter/reporter"
	"kubehostwarden/opscenter/user"
	"kubehostwarden/utils/middleware"
	"kubehostwarden/utils/responsor"
	"log"
	"net/http"
)

func NewServer() {
	var httpServer http.Server

	httpMux := http.NewServeMux()

	httpMux.HandleFunc("/health", health)
	httpMux.HandleFunc("/probe/register", probe.Register)

	httpMux.HandleFunc("/reporter/retrieve", reporter.Retrieve)

	httpMux.HandleFunc("/user/register", responsor.HandlePost(user.Register))
	httpMux.HandleFunc("/user/login", user.Login)

	httpMux.HandleFunc("/alarm/set", alarm.SetAlarm)

	httpServer.Handler = middleware.Cors(httpMux)
	httpServer.Addr = ":8080"

	fmt.Printf("Starting Http Server on port %s\n", httpServer.Addr)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
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
