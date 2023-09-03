package configservice

import (
	"encoding/json"
	"log"
	"mta2/modules/configservice/cinternals/constants"
	"mta2/modules/configservice/cpkg/ipconfig"
	"mta2/modules/natsmodule"
	"mta2/modules/utility"

	"github.com/nats-io/nats.go"
)

func PublishInvokeMessagetoNATS(c ipconfig.ConfigServiceIPMap, nc natsmodule.NATSConnInterface) error {

	_, err := nc.Subscribe(constants.INVOKE_SUB_SUBJECT_CONFIGSERVICE, HandlerMethodForPublishInvokeMessagetoNATS(c, nc))
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Started NATS Reply on Config Service")
	return nil
}

func HandlerMethodForPublishInvokeMessagetoNATS(c ipconfig.ConfigServiceIPMap, nc natsmodule.NATSConnInterface) nats.MsgHandler {
	return func(msg *nats.Msg) {
		log.Printf("Received message: %s\n", string(msg.Data))

		s := make([]*utility.Message, 0)

		for host, data := range c.GetValues() {
			s = append(s, &utility.Message{
				Hostname: host,
				Active:   data.ActiveIP,
			})
		}
		if encodedMessage, err := json.Marshal(s); err == nil {
			if err = nc.Publish(msg.Reply, encodedMessage); err != nil {
				log.Println("Error publishing NATS ", err)
			}
		} else {
			log.Println("Error Marshalling Invoke Data ", err)
		}
	}
}
