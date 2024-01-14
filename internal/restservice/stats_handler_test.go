package restservice

import (
	"context"
	"encoding/json"
	"github.com/MarouaneMan/github-api/api"
	"github.com/MarouaneMan/github-api/kvstore"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestStatsHandler(t *testing.T) {
	store := kvstore.NewInMemoryStore(kvstore.DefaultExpiration, kvstore.DefaultExpiration)

	_ = store.Write(context.Background(), "repositories", reposMock, kvstore.NoExpiration)

	req, err := http.NewRequest("GET", "/stats?language=golang", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := StatsHandler(store)
	err = handler(rr, req, nil)
	if err != nil {
		t.Fatalf("failed to handle stats request: %v", err)
		return
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	var statsResponse api.Stats
	err = json.Unmarshal(rr.Body.Bytes(), &statsResponse)
	if err != nil {
		t.Errorf("failed to unmarshal JSON response: %v", err)
	}

	expectedStats := api.Stats{
		Language:          "golang",
		TotalUsage:        3,
		TotalCodeSize:     3702,
		TotalRepositories: 4,
		AverageCodeSize:   1234,
	}

	if !reflect.DeepEqual(statsResponse, expectedStats) {
		t.Errorf("fetched stats do not match expected: got %+v, want %+v", statsResponse, expectedStats)
	}
}
