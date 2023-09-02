package utility

const (
	HOSTINGSERVICE       = "hostingservice"
	CONFIGSERVICE        = "configservice"
	HOSTINGSERVICEUSAAGE = ""
	CONFIGSERVICEUSAGE   = ""

	CONFIGSERVICE_PORT  = "CONFIGSERVICE_PORT"
	HOSTINGSERVICE_PORT = "HOSTINGSERVICE_PORT"
	NATS_URI            = "NATS_URI"
)

var (
	NATS_ADD string
	TaskChan = make(chan interface{}, 1)
)
