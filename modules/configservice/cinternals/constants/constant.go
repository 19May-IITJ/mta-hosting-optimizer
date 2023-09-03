package constants

import "sync"

// ENUMS for constant used in service as env variable and Default URL
const (
	DBPATH                           = "DBPATH"
	DEFAULTPATH                      = "/tmp/code/mock/data/ip2.json"
	Active                           = "1"
	Inactive                         = "0"
	Sep                              = "-"
	UPDATE_PUB_SUBJECT               = "hosting.update"
	CONFIGSERVICE_PUB_SUBJECT        = "config.down"
	INVOKE_SUB_SUBJECT_CONFIGSERVICE = "hosting.invoke"
)

var (
	Datamutex sync.Mutex
)
