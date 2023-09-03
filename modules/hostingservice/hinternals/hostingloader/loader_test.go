package hostingloader

import (
	"math/rand"
	"mta2/mock/mocking"
	"mta2/modules/hostingservice/hinternals/hostingconstants"
	"mta2/modules/hostingservice/pkg/dataconfig"
	"mta2/modules/utility"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoadConfigThreshold(t *testing.T) {
	t.Run("Test with default threshold value", func(t *testing.T) {
		os.Setenv(hostingconstants.MTA_THRESHOLD, "")
		defer os.Unsetenv(hostingconstants.MTA_THRESHOLD)
		result := LoadConfigThreshold()
		assert.Equal(t, 1, result)
	})
	t.Run("Test with random threshold value", func(t *testing.T) {
		testvalue := strconv.Itoa(rand.Int())
		os.Setenv(hostingconstants.MTA_THRESHOLD, testvalue)
		defer os.Unsetenv(hostingconstants.MTA_THRESHOLD)
		testthreshold, _ := strconv.Atoi(testvalue)
		result := LoadConfigThreshold()
		assert.Equal(t, testthreshold, result)
	})

}

func TestLoadActiveIPForHost(t *testing.T) {
	// Mock NATS connection and HostingServiceHostMap
	natsConn := new(mocking.MockNATSConn)
	hostMap := mocking.NewMockHostMap()
	expectedHostMap := dataconfig.NewHostMap()
	t.Run("Positive Test for LoadActiveIPForHost", func(t *testing.T) {

		expectedNATSResponse := &nats.Msg{
			Data: []byte(`[{"hostname":"mta-prod-1","active":2},{"hostname":"mta-prod-2","active":2},{"hostname":"mta-prod-3","active":1}]`),
		}
		// Set up expectations for your mocks
		natsConn.On("Request", hostingconstants.INVOKE_PUB_SUBJECT, []byte("Hello, Config Service!"), 10*time.Second).Return(expectedNATSResponse, nil)
		LoadActiveIPForHost(natsConn, hostMap, 10)

		expectedHostMap.Put(
			&utility.Message{
				Hostname: "mta-prod-1",
				Active:   2,
			},
			&utility.Message{
				Hostname: "mta-prod-2",
				Active:   2,
			},
			&utility.Message{
				Hostname: "mta-prod-3",
				Active:   1,
			})
		assert.Equal(t, expectedHostMap.GetValues(), hostMap.GetValues())
	})

}

func TestLoadUpdateStatusforHostName(t *testing.T) {
	mockNATSConn := new(mocking.MockNATSConn)
	mockHostMap := mocking.NewMockHostMap()
	// expectedHostMap := dataconfig.NewHostMap()

	// var v nats.MsgHandler
	// messageToRecieve := []byte(`{"hostname":"mta-prod-1","active":2}`)
	// v = func(msg *nats.Msg) {
	// 	// Simulate receiving a message
	// 	mockHostMap.On("Put", mock.AnythingOfType("*utility.Message")).Return()
	// 	// Invoke the callback function
	// 	message := utility.NewMessage()

	// 	if err := json.Unmarshal(msg.Data, &message); err == nil {
	// 		mockHostMap.Put(message)
	// 	}
	// }
	// v(&nats.Msg{
	// 	Data:    messageToRecieve,
	// 	Subject: hostingconstants.UPDATE_SUB_SUBJECT,
	// })
	t.Run("Positive Test for LoadUpdateStatusforHostName", func(t *testing.T) {

		mockNATSConn.On("Subscribe", hostingconstants.UPDATE_SUB_SUBJECT, mock.AnythingOfType("nats.MsgHandler")).
			Return(&nats.Subscription{}, nil)
		_, err := LoadUpdateStatusforHostName(mockNATSConn, mockHostMap)
		assert.Equal(t, err, nil)

	})
}
