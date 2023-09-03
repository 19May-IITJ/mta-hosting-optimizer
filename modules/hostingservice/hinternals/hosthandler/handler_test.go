package hosthandler

import (
	"context"
	"encoding/json"
	"mta2/mock/mocking"
	"mta2/modules/hostingservice/hinternals/hostingconstants"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveHostnames(t *testing.T) {
	// Mock NATS connection and HostingServiceHostMap
	natsConn := new(mocking.MockNATSConn)
	hostMap := mocking.NewMockHostMap()

	// Create a request and response recorder
	req, err := http.NewRequest("GET", "/hostnames", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Positive Test for RetrieveHostnames", func(t *testing.T) {

		w := httptest.NewRecorder()
		expectedResponse := &nats.Msg{
			Data: []byte(`[{"hostname":"mta-prod-1","active":2},{"hostname":"mta-prod-2","active":2},{"hostname":"mta-prod-3","active":1},{"hostname":"mta-prod-4","active":2}]`),
		}
		// Set up expectations for your mocks
		natsConn.On("Request", hostingconstants.INVOKE_PUB_SUBJECT, []byte("Hello, Config Service!"), 11*time.Second).Return(expectedResponse, nil)
		hostMap.On("IsEmpty").Return(true)

		// Create a context with a deadline
		ctx := context.Background()
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		req = req.WithContext(ctxWithTimeout)

		// Call the function with mocked dependencies
		handler := RetrieveHostnames(natsConn, 1, hostMap)
		handler(w, req)

		// Assert the response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse the response body
		var response []string
		err = json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		// Assert the response content
		assert.Equal(t, []string{"mta-prod-3"}, response)
	})
	t.Run("Positive Test for RetrieveHostnames with threshold 3", func(t *testing.T) {

		w := httptest.NewRecorder()
		expectedResponse := &nats.Msg{
			Data: []byte(`[{"hostname":"mta-prod-1","active":2},{"hostname":"mta-prod-2","active":2},{"hostname":"mta-prod-3","active":1},{"hostname":"mta-prod-4","active":2}]`),
		}
		// Set up expectations for your mocks
		natsConn.On("Request", hostingconstants.INVOKE_PUB_SUBJECT, []byte("Hello, Config Service!"), 11*time.Second).Return(expectedResponse, nil)
		hostMap.On("IsEmpty").Return(true)

		// Create a context with a deadline
		ctx := context.Background()
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		req = req.WithContext(ctxWithTimeout)

		// Call the function with mocked dependencies
		handler := RetrieveHostnames(natsConn, 0, hostMap)
		handler(w, req)

		// Assert the response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse the response body
		var response string
		err = json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		// Assert the response content
		assert.Equal(t, "No available Host have active MTA less than threshold 0", response)
	})
	t.Run("Negative Test for RetrieveHostnames", func(t *testing.T) {
		w := httptest.NewRecorder()
		expectedResponse := &nats.Msg{
			Data: []byte(`[{"hostname":"mta-prod-1","active":2},{"hostname":"mta-prod-2","active":2},{"hostname":"mta-prod-3","active":1},{"hostname":"mta-prod-4","active":2}]`),
		}
		// Set up expectations for your mocks
		natsConn.On("Request", hostingconstants.INVOKE_PUB_SUBJECT, []byte("Hello, Config Service!"), 11*time.Second).Return(expectedResponse, nil)
		hostMap.On("IsEmpty").Return(true)

		// Create a context with a deadline
		ctx := context.Background()
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 0*time.Second)
		defer cancel()
		req = req.WithContext(ctxWithTimeout)

		// Call the function with mocked dependencies
		handler := RetrieveHostnames(natsConn, 1, hostMap)
		handler(w, req)

		// Assert the response status code
		assert.Equal(t, http.StatusRequestTimeout, w.Code)

	})

	// Verify that your mocks were called as expected
	natsConn.AssertExpectations(t)
	hostMap.AssertExpectations(t)
}
