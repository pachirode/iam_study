package server

import (
	"net/http"

	"github.com/pachirode/iam_study/pkg/log"
)

func ServeHealthCheck(healthPath string, healthAddress string) {
	http.HandleFunc("/"+healthPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	if err := http.ListenAndServe(healthAddress, nil); err != nil {
		log.Fatalf("Error serving health check endpoint: %s", err.Error())
	}
}
