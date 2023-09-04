package handler

import (
	b "bytes"
	"context"
	"encoding/json"

	"mta2/mock/mocking"
	"mta2/modules/configservice/cinternals/constants"
	"mta2/modules/configservice/cinternals/loader"
	"mta2/modules/configservice/cpkg/ipconfig"
	"mta2/modules/utility"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// INTEGRATION TEST

// test to check possible positive & negative flow Refresh Data Set API
func TestRefreshDataSet(t *testing.T) {
	mocklist := ipconfig.NewIPConfigList()

	l := make([]*ipconfig.IPConfigData, 0)
	l = append(l, &ipconfig.IPConfigData{
		Hostname:    "dummy_1",
		IPAddresses: "127.0.0.1",
		Status:      true,
	}, &ipconfig.IPConfigData{
		Hostname:    "dummy_2",
		IPAddresses: "127.0.0.2",
		Status:      false,
	},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.3",
			Status:      false,
		},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.4",
			Status:      true,
		},
		&ipconfig.IPConfigData{
			Hostname:    "dummy_3",
			IPAddresses: "127.0.0.5",
			Status:      true,
		})
	mocklist.SetIPList(l)
	mockmap := ipconfig.NewMap()

	mockmap.Put("dummy_1", &ipconfig.HostData{
		HostedIP: []string{"127.0.0.1-1"},
		ActiveIP: 1,
	})
	mockmap.Put("dummy_2", &ipconfig.HostData{
		HostedIP: []string{"127.0.0.2-0"},
		ActiveIP: 0,
	})
	mockmap.Put("dummy_3", &ipconfig.HostData{
		HostedIP: []string{"127.0.0.3-0", "127.0.0.4-1", "127.0.0.5-1"},
		ActiveIP: 2,
	})
	loader.Ticker = time.NewTicker(30 * time.Second)
	natsConn := new(mocking.MockNATSConn)
	assert.Equal(t, 30, TTL)
	t.Run("Positive Test for Refersh Data Set", func(t *testing.T) {
		payload := []*ipconfig.IPConfigData{&ipconfig.IPConfigData{
			Hostname:    "dummy_1",
			IPAddresses: "127.0.0.1",
			Status:      false,
		},
		}

		bodyBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/refresh", b.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		bytes, _ := json.Marshal(&utility.Message{
			Hostname: "dummy_1",
			Active:   0,
		})
		natsConn.On("Publish", constants.UPDATE_PUB_SUBJECT, bytes).Return(nil)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist, natsConn))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK)
		mockConfig_expected := ipconfig.NewMap()
		mockIPList_expected := ipconfig.NewIPConfigList()

		mockConfig_expected.Put("dummy_1", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.1-0"},
			ActiveIP: 0,
		})
		mockConfig_expected.Put("dummy_2", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.2-0"},
			ActiveIP: 0,
		})
		mockConfig_expected.Put("dummy_3", &ipconfig.HostData{
			HostedIP: []string{"127.0.0.3-0", "127.0.0.4-1", "127.0.0.5-1"},
			ActiveIP: 2,
		})

		l := make([]*ipconfig.IPConfigData, 0)

		l = append(l, &ipconfig.IPConfigData{
			Hostname:    "dummy_1",
			IPAddresses: "127.0.0.1",
			Status:      false,
		}, &ipconfig.IPConfigData{
			Hostname:    "dummy_2",
			IPAddresses: "127.0.0.2",
			Status:      false,
		},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_3",
				IPAddresses: "127.0.0.3",
				Status:      false,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_3",
				IPAddresses: "127.0.0.4",
				Status:      true,
			},
			&ipconfig.IPConfigData{
				Hostname:    "dummy_3",
				IPAddresses: "127.0.0.5",
				Status:      true,
			})
		mockIPList_expected.SetIPList(l)
		assert.Equal(t, mockConfig_expected, mockmap)
	})

	t.Run("Timeout Test for Refersh Data Set", func(t *testing.T) {
		payload := []*ipconfig.IPConfigData{&ipconfig.IPConfigData{

			Hostname:    "dummy_1",
			IPAddresses: "127.0.0.1",
			Status:      true,
		},
		}
		bodyBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/refresh", b.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		bytes, _ := json.Marshal(&utility.Message{
			Hostname: "dummy_1",
			Active:   1,
		})
		natsConn.On("Publish", constants.UPDATE_PUB_SUBJECT, bytes).Return(nil)

		ctx, cancel := context.WithTimeout(req.Context(), 0*time.Microsecond)
		defer cancel()
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist, natsConn))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusRequestTimeout)
	})

	t.Run("IPs already marked for Refersh Data Set", func(t *testing.T) {
		payload := []*ipconfig.IPConfigData{&ipconfig.IPConfigData{

			Hostname:    "dummy_1",
			IPAddresses: "127.0.0.1",
			Status:      true,
		},
		}
		bodyBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/refresh", b.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		bytes, _ := json.Marshal(&utility.Message{
			Hostname: "dummy_1",
			Active:   1,
		})
		natsConn.On("Publish", constants.UPDATE_PUB_SUBJECT, bytes).Return(nil)

		ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist, natsConn))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("Invalid IP Hostname Mapping for Refersh Data Set", func(t *testing.T) {
		payload := []*ipconfig.IPConfigData{&ipconfig.IPConfigData{

			Hostname:    "dummy_2",
			IPAddresses: "127.0.0.1",
			Status:      true,
		},
		}
		bodyBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/refresh", b.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		bytes, _ := json.Marshal(&utility.Message{
			Hostname: "dummy_1",
			Active:   1,
		})
		natsConn.On("Publish", constants.UPDATE_PUB_SUBJECT, bytes).Return(nil)

		ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist, natsConn))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("Invalid Method for Refersh Data Set", func(t *testing.T) {
		payload := []*ipconfig.IPConfigData{&ipconfig.IPConfigData{

			Hostname:    "dummy_2",
			IPAddresses: "127.0.0.1",
			Status:      true,
		},
		}
		bodyBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("GET", "/refresh", b.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		bytes, _ := json.Marshal(&utility.Message{
			Hostname: "dummy_1",
			Active:   1,
		})
		natsConn.On("Publish", constants.UPDATE_PUB_SUBJECT, bytes).Return(nil)

		ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist, natsConn))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusMethodNotAllowed)
	})

	t.Run("Given Hostname Not Present for Refersh Data Set", func(t *testing.T) {
		payload := []*ipconfig.IPConfigData{&ipconfig.IPConfigData{

			Hostname:    "dummy_6",
			IPAddresses: "127.0.0.1",
			Status:      true,
		},
		}
		bodyBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/refresh", b.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		bytes, _ := json.Marshal(&utility.Message{
			Hostname: "dummy_1",
			Active:   1,
		})
		natsConn.On("Publish", constants.UPDATE_PUB_SUBJECT, bytes).Return(nil)

		ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist, natsConn))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest)
	})
	t.Run("Invalid Payload for Refersh Data Set", func(t *testing.T) {
		payload := []*ipconfig.IPConfigData{&ipconfig.IPConfigData{
			Status: true,
		},
		}
		bodyBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/refresh", b.NewReader(bodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		bytes, _ := json.Marshal(&utility.Message{
			Hostname: "dummy_1",
			Active:   1,
		})
		natsConn.On("Publish", constants.UPDATE_PUB_SUBJECT, bytes).Return(nil)

		ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist, natsConn))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest)
	})
}
