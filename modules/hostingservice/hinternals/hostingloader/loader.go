package hostingloader

import (
	"encoding/json"
	"log"
	"mta2/modules/hostingservice/hinternals/hostingconstants"
	"mta2/modules/hostingservice/pkg/dataconfig"
	"mta2/modules/utility"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	DEFAULT_TIMEOUT = 15
)

// Load Config Threshold loads the MTA_THRESHOLD env variable default:1
func LoadConfigThreshold() int {
	defaultThreshold := 1

	threshold := os.Getenv(hostingconstants.MTA_THRESHOLD)
	x, err := strconv.Atoi(threshold)
	if err != nil || x <= 0 {
		x = defaultThreshold
	}

	return x
}
func LoadActiveIPForHost(nc *nats.Conn, mp dataconfig.HostingServiceHostMap, timeout time.Duration) error {
	if timeout == 0 {
		timeout = DEFAULT_TIMEOUT
	}
	requestMsg := []byte("Hello, Config Service!")
	responseMsg, err := nc.Request(hostingconstants.INVOKE_PUB_SUBJECT, requestMsg, time.Second*timeout)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	if responseMsg != nil {
		result := make([]*utility.Message, 0)
		if err := json.Unmarshal(responseMsg.Data, &result); err == nil {
			mp.Put(result...)
		} else {
			log.Println("unable to se-serialze data ", err)
		}
	}
	return err
}

func LoadUpdateStatusforHostName(nc *nats.Conn, mp dataconfig.HostingServiceHostMap) {
	_, err := nc.Subscribe(hostingconstants.UPDATE_SUB_SUBJECT, func(msg *nats.Msg) {
		log.Printf("Received response: %s\n", string(msg.Data))
		dataconfig.DataMutex.Lock()
		message := utility.NewMessage()
		if err := json.Unmarshal(msg.Data, &message); err == nil {
			mp.Put(message)
		}
		dataconfig.DataMutex.Unlock()
	})
	if err != nil {
		log.Fatal(err)
	}
}
