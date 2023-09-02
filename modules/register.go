package modules

import (
	"log"
	"mta2/modules/configservice"
	"mta2/modules/configservice/cinternals/handler"
	"mta2/modules/configservice/cinternals/loader"
	"mta2/modules/hostingservice/hinternals/hostingloader"

	"mta2/modules/hostingservice/hinternals/hosthandler"
	"mta2/modules/hostingservice/pkg/dataconfig"

	"mta2/modules/configservice/cpkg/ipconfig"
	"mta2/modules/utility"
	"net"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

func RegisterService(serviceport string, kind string) {
	switch kind {
	case utility.CONFIGSERVICE:
		// ctx := context.Background()
		result := ipconfig.NewMap()
		list := ipconfig.NewIPConfigList()
		// Load configuration
		if err := loader.LoadConfigIPConfiguration(result, list); err == nil {
			// Create HTTP server
			log.Println("ip configurations loaded successfully")
			// register handlers to endpoints
			if nc, err := nats.Connect(utility.NATS_ADD); err == nil {
				configservice.PublishInvokeMessagetoNATS(result, nc)
				http.HandleFunc("/refresh", handler.RefreshDataSet(result, list, nc))
				log.Printf("Server listening on port %s\n", serviceport)
				// Starting server
				log.Fatal(http.ListenAndServe(hostAddr()+":"+serviceport, nil))
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
				hostingloader.LoadUpdateStatusforHostName(nc, mp)
				log.Println("ip configurations loaded successfully")
				// register handlers to endpoints
				http.HandleFunc("/hostnames", hosthandler.RetrieveHostnames(nc, x, mp))
				log.Printf("Server listening on port %s\n", serviceport)
				// Starting server
				log.Fatal(http.ListenAndServe(hostAddr()+":"+serviceport, nil))
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
