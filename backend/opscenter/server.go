package opscenter

import (
	"kubehostwarden/opscenter/alarm"
	opsHost "kubehostwarden/opscenter/host"
	"kubehostwarden/opscenter/logger"
	"kubehostwarden/opscenter/reporter"
	"kubehostwarden/opscenter/user"
	log1 "kubehostwarden/utils/log"
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
	// alarm api
	authMux.HandleFunc("/alarm/setthreshold", responsor.HandlePost(alarm.SetThreshold))
	authMux.HandleFunc("/alarm/deletethreshold", responsor.HandlePost(alarm.DeleteThreshold))
	authMux.HandleFunc("/alarm/getthreshold", responsor.HandleGet(alarm.GetThreshold))
	// logger api
	authMux.HandleFunc("/logger/get", responsor.HandleGet(logger.GetLogs))

	authHandler := middleware.Auth(authMux)

	mainMux.HandleFunc("/user/register", responsor.HandlePost(user.Register))
	mainMux.HandleFunc("/user/login", responsor.HandlePost(user.Login))

	mainMux.HandleFunc("/reporter/report", responsor.HandleGet(reporter.Report))
	mainMux.HandleFunc("/health", health)

	// httpNoAuthMux.HandleFunc("/alarm/set", alarm.SetAlarm)

	mainMux.Handle("/", authHandler)

	httpServer.Handler = middleware.Cors(mainMux)
	httpServer.Addr = ":8080"

	log1.Info("Opscenter server started", "addr", httpServer.Addr)
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
