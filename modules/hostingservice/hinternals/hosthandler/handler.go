package hosthandler

import (
	"context"
	"encoding/json"
	"log"
	"mta2/modules/hostingservice/hinternals/hostingloader"
	"mta2/modules/hostingservice/pkg/dataconfig"

	"net/http"
	"time"

	"github.com/nats-io/nats.go"
)

var counter int

// Retrieve Hostnames handler return http handleFunc used to get inefficient hostnames having active no. of IP <= threshold value
func RetrieveHostnames(nc *nats.Conn, maxIPs int, result dataconfig.HostingServiceHostMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusRequestTimeout)
			return

		default:
			if result.IsEmpty() {
				counter++
				hostingloader.LoadActiveIPForHost(nc, result)
				if result.IsEmpty() && counter > 5 {
					http.Error(w, "seems config service is down", http.StatusInternalServerError)
					return
				} else {
					http.Error(w, "can't serve request no data available", http.StatusInternalServerError)

				}
			} else {
				log.Println("received request Retrieve Hostnames")
				inefficientHostnames := dataconfig.GetHostnamesWithMaxIPs(maxIPs, result)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(inefficientHostnames)
			}
		}
	}
}
