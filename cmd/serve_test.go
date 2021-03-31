package cmd

import (
	"github.com/rzauls/curfetch/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_healthHandler(t *testing.T) {
	storage := db.NewMockStorage()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	t.Run("health handler successfully responds", func(t *testing.T) {
		s := &Server{db: storage}
		s.healthHandler(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
		}
	})

}