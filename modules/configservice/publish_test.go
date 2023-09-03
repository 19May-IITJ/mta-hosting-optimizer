package configservice

import (
	"encoding/json"
	"mta2/mock/mocking"
	"mta2/modules/configservice/cinternals/constants"
	"mta2/modules/configservice/cpkg/ipconfig"
	"mta2/modules/utility"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPublishInvokeMessagetoNATS(t *testing.T) {
	mockMap := ipconfig.NewMap()
	natsConn := new(mocking.MockNATSConn)
	natsConn.On("Subscribe", constants.INVOKE_SUB_SUBJECT_CONFIGSERVICE, mock.AnythingOfType("nats.MsgHandler")).Return(&nats.Subscription{}, nil)
	err := PublishInvokeMessagetoNATS(mockMap, natsConn)
	assert.NoError(t, err)
}

func TestHandlerMethodForPublishInvokeMessagetoNATS(t *testing.T) {
	s := make([]*utility.Message, 0)
	mockMap := ipconfig.NewMap()
	natsConn := new(mocking.MockNATSConn)
	bytes, _ := json.Marshal(s)

	v := nats.Msg{
		Reply: "dummy.topic",
	}
	natsConn.On("Publish", v.Reply, bytes).Return(nil)
	handler := HandlerMethodForPublishInvokeMessagetoNATS(mockMap, natsConn)
	handler(&v)
}
