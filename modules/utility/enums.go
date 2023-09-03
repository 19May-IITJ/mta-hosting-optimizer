package utility

const (
	HOSTINGSERVICE       = "hostingservice"
	CONFIGSERVICE        = "configservice"
	HOSTINGSERVICEUSAAGE = "hosts API handler /hostname for uncovering in-efficient server as per MTA_THRESHOLD value"
	CONFIGSERVICEUSAGE   = "host as a DB interface layer and do hosts API handler provide /refresh to change data in local and DB(JSON file)"

	CONFIGSERVICE_PORT  = "CONFIGSERVICE_PORT"
	HOSTINGSERVICE_PORT = "HOSTINGSERVICE_PORT"
	NATS_URI            = "NATS_URI"
)

var (
	NATS_ADD string
	TaskChan = make(chan interface{}, 1)
)
