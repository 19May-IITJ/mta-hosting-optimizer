package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"mta2/modules/configservice/cinternals/constants"
	"mta2/modules/configservice/cinternals/loader"
	"mta2/modules/natsmodule"

	"mta2/modules/configservice/cpkg/ipconfig"
	"mta2/modules/utility"
	"net/http"
	"strings"

	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

const (
	TTL = 30
)

// Refresh DataSet handler return http handleFunc used to reload all ip & hostname data and active ip's under hostname
func RefreshDataSet(c ipconfig.ConfigServiceIPMap, ipl ipconfig.IPListI, nc natsmodule.NATSConnInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

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
		validate := validator.New()
		validate.RegisterValidation("string", ipconfig.IsString)
		for i, ipConfig := range requestBody {
			if err := validate.Struct(ipConfig); err != nil {
				log.Printf("Validation error for element %d and value %v: %v\n", i, ipConfig, err)
				http.Error(w, "Validation error for element", http.StatusBadRequest)
				return
			}

		}

		if len(requestBody) > 0 {

			for _, entry := range requestBody {
				if c.Contains(entry.Hostname) {
					hd, _ := c.GetValue(entry.Hostname)
					var s string
					var o string
					if entry.Status {
						s = strings.Join([]string{entry.IPAddresses, constants.Active}, constants.Sep)
						o = strings.Join([]string{entry.IPAddresses, constants.Inactive}, constants.Sep)
					} else {
						s = strings.Join([]string{entry.IPAddresses, constants.Inactive}, constants.Sep)
						o = strings.Join([]string{entry.IPAddresses, constants.Active}, constants.Sep)

					}
					if !slices.Contains(hd.HostedIP, s) {
						if index := slices.Index(hd.HostedIP, o); index != -1 {
							if entry.Status {
								hd.ActiveIP++
							} else {
								hd.ActiveIP--
							}
							hd.HostedIP[index] = s
							go func(update *ipconfig.IPConfigData) {
								constants.Datamutex.Lock()
								loader.FLAGTOSAVE = true
								loader.Ticker.Reset(TTL * time.Second)
								if index := loader.Search(ipl.GetIPValues(), update.IPAddresses); index != -1 {
									ipl.GetIPValues()[index] = &ipconfig.IPConfigData{
										Hostname:    update.Hostname,
										IPAddresses: update.IPAddresses,
										Status:      update.Status,
									}
								}
								log.Printf("Added %s status %v success for Host %s\n", update.IPAddresses, update.Status, update.Hostname)
								constants.Datamutex.Unlock()
							}(entry)
							if message, err := json.Marshal(&utility.Message{
								Hostname: entry.Hostname,
								Active:   hd.ActiveIP,
							}); err == nil {
								if err = nc.Publish(constants.UPDATE_PUB_SUBJECT, message); err != nil {
									log.Println("Error Publishing message to NATS ", err)
								}
								log.Println("published message success ", string(message))
							} else {
								log.Println("Unable to parse data couldn't publish to NATS ", err)
							}
						} else {
							http.Error(w, "invalid IP and Hostname Mapping", http.StatusBadRequest)
							return
						}
					} else {
						http.Error(w, "IP already have marked with the status provided", http.StatusBadRequest)
						return
					}
				} else {
					http.Error(w, "Given Hostname not present in Data Base", http.StatusBadRequest)
					return
				}
			}
		} else {
			http.Error(w, "Unable to Refresh Data no Data provided to Refresh", http.StatusExpectationFailed)
			return
		}

		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusRequestTimeout)
			return
		default:
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Data Updated successfully"))
		}
	}
}
