package opscenter

import (
	opsHost "kubehostwarden/opscenter/host"
	"kubehostwarden/opscenter/reporter"
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
	// host api
	authMux.HandleFunc("/host/register", responsor.HandlePost(opsHost.Register))
	authMux.HandleFunc("/host/delete", responsor.HandlePost(opsHost.Delete))
	authMux.HandleFunc("/host/retrieve", responsor.HandleGet(opsHost.Retrieve))
	// user api
	authMux.HandleFunc("/user/retrieve", responsor.HandleGet(user.Retrieve))
	// reporter api
	authMux.HandleFunc("/reporter/report", reporter.Report)

	authHandler := middleware.Auth(authMux)

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
