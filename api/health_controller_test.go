package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gosom/goappbuild/api"
	"github.com/stretchr/testify/require"
)

func Test_HealthController_GetHealth(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req := getHTTPRequest(ctx, t, http.MethodGet, "/health", nil)

	rr := httptest.NewRecorder()

	hc := api.NewHealthController()

	handler := http.HandlerFunc(hc.GetHealth)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp api.HealthResponse

	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	require.Equal(t, "ok", resp.Status)
}
