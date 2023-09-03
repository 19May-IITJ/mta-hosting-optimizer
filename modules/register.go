package modules

import (
	"context"
	"fmt"
	"log"
	"mta2/modules/configservice"
	"mta2/modules/configservice/cinternals/handler"
	"mta2/modules/configservice/cinternals/loader"
	"mta2/modules/hostingservice/hinternals/hostingloader"
	"time"

	"mta2/modules/hostingservice/hinternals/hosthandler"
	"mta2/modules/hostingservice/pkg/dataconfig"

	"mta2/modules/configservice/cpkg/ipconfig"
	"mta2/modules/utility"
	"net"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

func RegisterService(ctx context.Context, serviceport string, kind string, s *http.Server) {
	switch kind {
	case utility.CONFIGSERVICE:
		result := ipconfig.NewMap()
		list := ipconfig.NewIPConfigList()
		// Load configuration
		if err := loader.LoadConfigIPConfiguration(result, list); err == nil {
			// Create HTTP server
			log.Println("ip configurations loaded successfully")
			// register handlers to endpoints
			if nc, err := nats.Connect(utility.NATS_ADD); err == nil {
				loader.Ticker = time.NewTicker(handler.TTL * time.Second)
				go loader.TTLForFileSaving(ctx, list)
				configservice.PublishInvokeMessagetoNATS(result, nc)
				http.HandleFunc("/refresh", handler.RefreshDataSet(result, list, nc))
				log.Printf("Server listening on port %s\n", serviceport)
				// Starting server
				s.Addr = hostAddr() + ":" + serviceport

				go func() {
					// Start the server and listen for incoming requests.
					if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						fmt.Printf("Error: %v\n", err)
					}
				}()

			} else {
				log.Printf("Error %v \n-*-unable to launch application-*-\n", err)
			}
		} else {
			log.Printf("Error %v \n-*-unable to launch application-*-\n", err)
		}
	case utility.HOSTINGSERVICE:

		mp := dataconfig.NewHostMap()
		// Load configuration
		x := hostingloader.LoadConfigThreshold()
		if nc, err := nats.Connect(utility.NATS_ADD); err == nil {
			if err := hostingloader.LoadActiveIPForHost(nc, mp, 0); err == nil || err == nats.ErrTimeout {
				// Create HTTP server
				if _, err := hostingloader.LoadUpdateStatusforHostName(nc, mp); err != nil {
					log.Fatalf("Error %v \n-*-unable to launch application-*-\n", err)
				}
				log.Println("ip configurations loaded successfully")
				// register handlers to endpoints
				http.HandleFunc("/hostnames", hosthandler.RetrieveHostnames(nc, x, mp))
				log.Printf("Server listening on port %s\n", serviceport)
				// Starting server
				s.Addr = hostAddr() + ":" + serviceport

				go func() {
					// Start the server and listen for incoming requests.
					if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						fmt.Printf("Error: %v\n", err)
					}
				}()
			} else {
				log.Printf("Error %v \n-*-unable to launch application-*-\n", err)
			}
		}
	}
}
func hostAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
