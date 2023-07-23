package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func getHTTPRequest(ctx context.Context, tb testing.TB, method, u string, payload any) *http.Request {
	tb.Helper()

	var body io.Reader
	if payload == nil {
		body = http.NoBody
	} else {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(payload)
		require.NoError(tb, err)

		body = &buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u, body)
	require.NoError(tb, err)

	return req
}
