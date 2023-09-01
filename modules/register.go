package modules

import (
	"log"
	"mta2/internal/config"
	"mta2/internal/handler"
	"mta2/modules/utility"
	"mta2/pkg/ipconfig"
	"net"
	"net/http"
	"os"
)

func RegisterService(serviceport string, kind string) {
	switch kind {
	case utility.CONFIGSERVICE:

		result := ipconfig.NewMap()
		list := ipconfig.NewIPConfigList()
		// Load configuration
		x := config.LoadConfigThreshold()
		if err := config.LoadConfigIPConfiguration(result, list); err == nil {
			// Create HTTP server
			log.Println("ip configurations loaded successfully")
			// register handlers to endpoints
			http.HandleFunc("/hostnames", handler.RetrieveHostnames(x, result))
			http.HandleFunc("/refresh", handler.RefreshDataSet(result, list))
			log.Printf("Server listening on port %s\n", serviceport)
			// Starting server
			log.Fatal(http.ListenAndServe(hostAddr()+":"+serviceport, nil))
		} else {
			log.Printf("Error %v \n-*-unable to launch application-*-\n", err)
		}
	case utility.HOSTINGSERVICE:
		result := ipconfig.NewMap()
		list := ipconfig.NewIPConfigList()
		// Load configuration
		x := config.LoadConfigThreshold()
		if err := config.LoadConfigIPConfiguration(result, list); err == nil {
			// Create HTTP server
			log.Println("ip configurations loaded successfully")
			// register handlers to endpoints
			http.HandleFunc("/hostnames", handler.RetrieveHostnames(x, result))
			http.HandleFunc("/refresh", handler.RefreshDataSet(result, list))
			log.Printf("Server listening on port %s\n", serviceport)
			// Starting server
			log.Fatal(http.ListenAndServe(hostAddr()+":"+serviceport, nil))
		} else {
			log.Printf("Error %v \n-*-unable to launch application-*-\n", err)
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
