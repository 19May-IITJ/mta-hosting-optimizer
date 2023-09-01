package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mta2/internal/loader"
	"mta2/pkg/ipconfig"
	"net/http"
	"time"
)

// Retrieve Hostnames handler return http handleFunc used to get inefficient hostnames having active no. of IP <= threshold value
func RetrieveHostnames(maxIPs int, result ipconfig.ConfigServiceIPMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusRequestTimeout)
			return

		default:
			log.Println("received request Retrieve Hostnames")
			inefficientHostnames := ipconfig.GetHostnamesWithMaxIPs(maxIPs, result)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(inefficientHostnames)
		}
	}
}

// Refresh DataSet handler return http handleFunc used to reload all ip & hostname data and active ip's under hostname
func RefreshDataSet(c ipconfig.ConfigServiceIPMap, ipl ipconfig.IPListI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusRequestTimeout)
			return
		default:
			log.Println("received request Refresh Data Set")

			if r.Method != http.MethodPost {
				http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
				return
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			requestBody := []*ipconfig.IPConfigData{}
			if err := json.Unmarshal(body, &requestBody); err != nil {
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				return
			}
			if len(requestBody) > 0 {

				for _, entry := range requestBody {
					if c.Contains(entry.IPAddresses) {
						state, _ := c.GetValue(entry.IPAddresses)
						if state.Hostname == entry.Hostname {
							c.Put(entry.IPAddresses, &ipconfig.IPState{
								State:    entry.Status,
								Hostname: entry.Hostname,
							})

							if index := loader.Search(ipl.GetIPValues(), entry.IPAddresses); index != -1 {
								ipl.GetIPValues()[index] = &ipconfig.IPConfigData{
									Hostname:    entry.Hostname,
									IPAddresses: entry.IPAddresses,
									Status:      entry.Status,
								}
							}
						} else {
							http.Error(w, "Invalid IP and Hostname Mapping", http.StatusBadRequest)
							return
						}
					} else {
						http.Error(w, "Given IP not present in Data Base", http.StatusBadRequest)
						return
					}
				}
			} else {
				statusCode := http.StatusExpectationFailed
				w.WriteHeader(statusCode)
				responseBody := fmt.Sprintf("Unable to Refresh Data status code: %d, Error %s", statusCode, err.Error())
				w.Write([]byte(responseBody))
			}
		}
	}
}
