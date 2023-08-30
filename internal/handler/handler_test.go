package handler

import (
	"encoding/json"
	"mta2/internal/constants"
	"mta2/pkg/ipconfig"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// INTEGRATION TEST

// test to check possible positive & negative flow Refresh Data Set API
func TestRefreshDataSet(t *testing.T) {
	mocklist := ipconfig.NewIPConfigList()
	req, err := http.NewRequest("GET", "/refresh", nil)
	if err != nil {
		t.Fatal(err)
	}
	l := make([]*ipconfig.IPConfig, 0)
	l = append(l, &ipconfig.IPConfig{
		Hostname:    "dummy_1",
		IPAddresses: "127.0.0.1",
		Status:      true,
	}, &ipconfig.IPConfig{
		Hostname:    "dummy_2",
		IPAddresses: "127.0.0.1",
		Status:      false,
	})
	mocklist.SetIPList(l)
	mockmap := ipconfig.NewMap()
	mockmap.Put("dummy_1", 1)
	mockmap.Put("dummy_2", 2)
	mockmap.Put("dummy_3", 0)

	t.Run("Positive Test for Refersh Data Set", func(t *testing.T) {
		os.Setenv(constants.DBPATH, "/Users/b0268986/mta2/mock/test_data/ipconfig_test.json")
		defer os.Unsetenv(constants.DBPATH)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK)
		mockConfig_expected := ipconfig.NewMap()
		mockIPList_expected := ipconfig.NewIPConfigList()
		mockConfig_expected.Put("dummy_1", 1)
		mockConfig_expected.Put("dummy_2", 2)
		mockConfig_expected.Put("dummy_3", 0)
		l := make([]*ipconfig.IPConfig, 0)
		l = append(l,
			&ipconfig.IPConfig{
				Hostname:    "dummy_1",
				IPAddresses: "127.0.0.1",
				Status:      false,
			},
			&ipconfig.IPConfig{
				Hostname:    "dummy_1",
				IPAddresses: "127.0.0.2",
				Status:      false,
			},
			&ipconfig.IPConfig{
				Hostname:    "dummy_1",
				IPAddresses: "127.0.0.3",
				Status:      true,
			},
			&ipconfig.IPConfig{
				Hostname:    "dummy_2",
				IPAddresses: "127.0.0.4",
				Status:      true,
			},
			&ipconfig.IPConfig{
				Hostname:    "dummy_2",
				IPAddresses: "127.0.0.5",
				Status:      true,
			},
			&ipconfig.IPConfig{
				Hostname:    "dummy_3",
				IPAddresses: "127.0.0.6",
				Status:      false,
			},
		)
		mockIPList_expected.SetIPList(l)
		assert.Equal(t, mockConfig_expected, mockmap)
		assert.Equal(t, mockIPList_expected, mocklist)
		// if status := rr.Code; status != http.StatusOK {
		// 	t.Errorf("handler returned wrong status code: got %v want %v",
		// 		status, http.StatusOK)
		// }
	})
	t.Run("Negative Test for Refersh Data Set", func(t *testing.T) {
		rr := httptest.NewRecorder()
		os.Setenv(constants.DBPATH, "random/random")
		defer os.Unsetenv(constants.DBPATH)
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusExpectationFailed)
	})
}

// test to check possible positive & negative flow Retrieve Hostnames API

func TestRetrieveHostnames(t *testing.T) {
	/* In Retrieve data we expect that data is loaded either
	via default service up or via refresh API
	*/
	req, err := http.NewRequest("GET", "/hostnames", nil)
	if err != nil {
		t.Fatal(err)
	}
	mockmap := ipconfig.NewMap()
	threshold := 1
	mockmap.Put("dummy_1", 1)
	mockmap.Put("dummy_2", 2)
	mockmap.Put("dummy_3", 0)
	inefficientHostname_expected := make([]string, 0)
	inefficientHostname_expected = append(inefficientHostname_expected, "dummy_1", "dummy_3")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RetrieveHostnames(threshold, mockmap))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	var response []string
	json.Unmarshal(rr.Body.Bytes(), &response)
	sort.Slice(response, func(p, q int) bool {
		return response[p] < response[q]
	})
	assert.Equal(t, inefficientHostname_expected, response)
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }
}
