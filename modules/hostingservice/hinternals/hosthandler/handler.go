package hosthandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mta2/modules/hostingservice/hinternals/hostingloader"
	"mta2/modules/hostingservice/pkg/dataconfig"
	"mta2/modules/natsmodule"

	"net/http"
	"time"
)

var counter int

const (
	dEFAULTCONTEXTTIMEOUT = 10
)

// Retrieve Hostnames handler return http handleFunc used to get inefficient hostnames having active no. of IP <= threshold value
func RetrieveHostnames(nc natsmodule.NATSConnInterface, maxIPs int, result dataconfig.HostingServiceHostMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx, cancel := context.WithTimeout(r.Context(), dEFAULTCONTEXTTIMEOUT*time.Second)
		defer cancel()
		log.Println("received request Retrieve Hostnames")
		var inefficientHostnames []string
		if result.IsEmpty() {
			counter++
			err = hostingloader.LoadActiveIPForHost(nc, result, dEFAULTCONTEXTTIMEOUT+1)
		}
		inefficientHostnames = dataconfig.GetHostnamesWithMaxIPs(maxIPs, result)

		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				if result.IsEmpty() && counter > 3 {
					http.Error(w, "seems config service is down", http.StatusInternalServerError)
				} else {
					http.Error(w, "Request timed out no data available", http.StatusRequestTimeout)
				}
				return
			}
		default:
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			if len(inefficientHostnames) > 0 {
				json.NewEncoder(w).Encode(inefficientHostnames)
			} else {
				var response string
				if err == nil {
					response = fmt.Sprintf("No available Host have active MTA less than threshold %v", maxIPs)
				} else {
					response = fmt.Sprintf("seems config service got down and no NATS sub avialable %v", err)
				}
				json.NewEncoder(w).Encode(response)
			}
		}
	}
}
