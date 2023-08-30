package handler

import (
	"mta2/internal/constants"
	"mta2/pkg/ipconfig"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshDataSet(t *testing.T) {
	mocklist := ipconfig.NewIPConfigList()
	req, err := http.NewRequest("GET", "/refresh", nil)
	if err != nil {
		t.Fatal(err)
	}
	l := make([]*ipconfig.IPConfig, 0)
	l = append(l, &ipconfig.IPConfig{
		Hostname:    "dummy1",
		IPAddresses: "127.0.0.1",
		Status:      true,
	}, &ipconfig.IPConfig{
		Hostname:    "dummy2",
		IPAddresses: "127.0.0.1",
		Status:      false,
	})
	mocklist.SetIPList(l)
	mockmap := ipconfig.NewMap()
	mockmap.Put("dummy_1", 1)
	mockmap.Put("dummy_2", 2)
	mockmap.Put("dummy_3", 0)

	t.Run("Positive Test for Refersh Data Set", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RefreshDataSet(mockmap, mocklist))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK)
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
func TestRetrieveHostnames(t *testing.T) {
	req, err := http.NewRequest("GET", "/hostnames", nil)
	if err != nil {
		t.Fatal(err)
	}
	mockmap := ipconfig.NewMap()
	threshold := 1
	mockmap.Put("dummy_1", 1)
	mockmap.Put("dummy_2", 2)
	mockmap.Put("dummy_3", 0)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RetrieveHostnames(threshold, mockmap))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
