package constants

import "sync"

// ENUMS for constant used in service as env variable and Default URL
const (
	DBPATH                           = "DBPATH"
	DEFAULTPATH                      = "/Users/b0268986/mta2/mock/data/ipconfig.json"
	Active                           = "1"
	Inactive                         = "0"
	Sep                              = "-"
	UPDATE_PUB_SUBJECT               = "hosting.update"
	INVOKE_SUB_SUBJECT_CONFIGSERVICE = "hosting.invoke"
)

var (
	Datamutex sync.Mutex
)
