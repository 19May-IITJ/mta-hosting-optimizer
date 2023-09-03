package hostingloader

import (
	"encoding/json"
	"log"
	"mta2/modules/hostingservice/hinternals/hostingconstants"
	"mta2/modules/hostingservice/pkg/dataconfig"
	"mta2/modules/natsmodule"
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

/*
"Sync" NATS Request Reply Engine to hit Config at time of service registry
and API hit of getting /hostnames provided the local cache is empty.
It panics if the atleast one NATS subsciption either from Config Service end or NATS cli is not present
*/
func LoadActiveIPForHost(nc natsmodule.NATSConnInterface, mp dataconfig.HostingServiceHostMap, timeout time.Duration) error {
	if timeout == 0 {
		timeout = DEFAULT_TIMEOUT
	}
	requestMsg := []byte("Hello, Config Service!")
	responseMsg, err := nc.Request(hostingconstants.INVOKE_PUB_SUBJECT, requestMsg, time.Second*timeout)
	if responseMsg != nil {
		result := make([]*utility.Message, 0)
		if err := json.Unmarshal(responseMsg.Data, &result); err == nil {
			dataconfig.DataMutex.Lock()
			mp.Put(result...)
			dataconfig.DataMutex.Unlock()
		} else {
			log.Println("unable to se-serialze data ", err)
		}
	}
	return err
}

// "Async" NATS subscribition for listening to update signal for Config Service caused via its Refresh Data API hit
func LoadUpdateStatusforHostName(nc natsmodule.NATSConnInterface, mp dataconfig.HostingServiceHostMap) (*nats.Subscription, error) {
	return nc.Subscribe(hostingconstants.UPDATE_SUB_SUBJECT, HandlerForLoadUpdateStatusforHostName(mp))
}

// "Async" NATS subscribition for listening to down signal of Config Service
func RollBackDataONConfigDown(nc natsmodule.NATSConnInterface, mp dataconfig.HostingServiceHostMap) (*nats.Subscription, error) {
	return nc.Subscribe(hostingconstants.HOSTINGCONFIG_SUB_SUBJECT, HandlerForRollBackDataONConfigDown(mp))
}

func HandlerForRollBackDataONConfigDown(mp dataconfig.HostingServiceHostMap) nats.MsgHandler {
	return func(msg *nats.Msg) {
		log.Printf("Received response: %s\n", string(msg.Data))
		dataconfig.DataMutex.Lock()
		mp.Clear()
		dataconfig.DataMutex.Unlock()
	}
}

func HandlerForLoadUpdateStatusforHostName(mp dataconfig.HostingServiceHostMap) nats.MsgHandler {
	return func(msg *nats.Msg) {
		log.Printf("Received response: %s\n", string(msg.Data))
		dataconfig.DataMutex.Lock()
		message := utility.NewMessage()
		if err := json.Unmarshal(msg.Data, &message); err == nil {
			mp.Put(message)
		}
		dataconfig.DataMutex.Unlock()
	}
}
