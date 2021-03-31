package cmd

import (
	"github.com/gorilla/mux"
	"github.com/rzauls/curfetch/db"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type args struct {
	w http.ResponseWriter
	r *http.Request
}
func TestServer_healthHandler(t *testing.T) {
	storage := db.NewMockStorage([]db.Currency{}, time.Now())
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	t.Run("successfully responds", func(t *testing.T) {
		s := &Server{db: storage}
		s.healthHandler(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("got status %d but wanted %d", w.Code, http.StatusOK)
		}
	})

}

func TestServer_newestHandler(t *testing.T) {
	emptyStorage := db.NewMockStorage([]db.Currency{}, time.Now())
	timestamp := time.Now()
	dataPoints := []db.Currency{
		{
			Code:    "USD",
			Value:   "1.11",
			PubDate: timestamp,
		},
		{
			Code:    "AUD",
			Value:   "2.22",
			PubDate: timestamp,
		},
		{
			Code:    "EUR",
			Value:   "3.33",
			PubDate: timestamp,
		},
	}
	validStorage := db.NewMockStorage(dataPoints, timestamp)

	tests := []struct {
		name         string
		storage      db.Storage
		responseCode int
	}{
		{
			name:         "200 success",
			storage:      validStorage,
			responseCode: http.StatusOK,
		},
		{
			name:         "204 with no data",
			storage:      emptyStorage,
			responseCode: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/newest", nil)
			s := &Server{db: tt.storage}
			s.newestHandler(w, r)
			if w.Code != tt.responseCode {
				t.Errorf("got status %d but wanted %d", w.Code, tt.responseCode)
			}
		})
	}
}

func TestServer_historyHandler(t *testing.T) {
	emptyStorage := db.NewMockStorage([]db.Currency{}, time.Now())
	
	timestamp := time.Now()
	dataPoints := []db.Currency{
		{
			Code:    "USD",
			Value:   "1.11",
			PubDate: time.Now().Add(-2 * time.Minute),
		},
		{
			Code:    "USD",
			Value:   "2.22",
			PubDate: time.Now().Add(-1 * time.Minute),
		},
		{
			Code:    "USD",
			Value:   "3.33",
			PubDate: time.Now(),
		},
	}
	validStorage := db.NewMockStorage(dataPoints, timestamp)
	
	tests := []struct {
		name         string
		storage      db.Storage
		responseCode int
		currencyCode string
	}{
		{
			name:         "returns correct data",
			storage:      validStorage,
			responseCode: http.StatusOK,
			currencyCode: "usd",
		},
		{
			name:         "204 with no data",
			storage:      emptyStorage,
			responseCode: http.StatusNoContent,
			currencyCode: "usd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				db: tt.storage,
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/history", nil)
			r = mux.SetURLVars(r, map[string]string{
				"currency": tt.currencyCode,
			})

			s.historyHandler(w, r)
			if w.Code != tt.responseCode {
				t.Errorf("got status %d but wanted %d", w.Code, tt.responseCode)
			}
			// TODO: test if response contains correct currency responseCode
		})
	}
}