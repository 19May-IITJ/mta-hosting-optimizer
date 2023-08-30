package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"mta2/internal/config"
	"mta2/pkg/ipconfig"
	"net/http"
)

// Retrieve Hostnames handler return http handleFunc used to get inefficient hostnames having active no. of IP <= threshold value
func RetrieveHostnames(maxIPs int, result ipconfig.Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("received request Retrieve Hostnames")
		inefficientHostnames := ipconfig.GetHostnamesWithMaxIPs(maxIPs, result)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inefficientHostnames)
	}
}

// Refresh DataSet handler return http handleFunc used to reload all ip & hostname data and active ip's under hostname
func RefreshDataSet(c ipconfig.Configuration, ipl ipconfig.IPListI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("received request Refresh Data Set")
		if !c.IsEmpty() {
			c.Clear()
		}
		if !ipl.IsEmpty() {
			ipl.Clear()
		}
		if err := config.LoadConfigIPConfiguration(c, ipl); err == nil {
			statusCode := http.StatusOK
			w.WriteHeader(statusCode)
			responseBody := fmt.Sprintf("Data Refreshed Successfully status code: %d", statusCode)
			w.Write([]byte(responseBody))

		} else {
			statusCode := http.StatusExpectationFailed
			w.WriteHeader(statusCode)
			responseBody := fmt.Sprintf("Unable to Refresh Data status code: %d, Error %s", statusCode, err.Error())
			w.Write([]byte(responseBody))
		}
	}
}
