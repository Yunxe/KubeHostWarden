package opscenter

import (
	"kubehostwarden/opscenter/probe"
	"kubehostwarden/opscenter/user"
	"kubehostwarden/utils/logger"
	"kubehostwarden/utils/middleware"
	"kubehostwarden/utils/responsor"
	"log"
	"net/http"
)

func NewServer() {
	var httpServer http.Server
	mainMux := http.NewServeMux()
	authMux := http.NewServeMux()

	authMux.HandleFunc("/probe/register", responsor.HandlePost(probe.Register))

	authHandler := middleware.Auth(authMux)

	// httpNoAuthMux.HandleFunc("/health", health)
	// httpNoAuthMux.HandleFunc("/probe/register", probe.Register)

	// httpNoAuthMux.HandleFunc("/reporter/retrieve", reporter.Retrieve)

	mainMux.HandleFunc("/user/register", responsor.HandlePost(user.Register))
	mainMux.HandleFunc("/user/login", responsor.HandlePost(user.Login))
	mainMux.HandleFunc("/health", health)

	// httpNoAuthMux.HandleFunc("/alarm/set", alarm.SetAlarm)

	mainMux.Handle("/", authHandler)

	httpServer.Handler = middleware.Cors(mainMux)
	httpServer.Addr = ":8080"

	logger.Info("Opscenter server started", "addr", httpServer.Addr)
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
