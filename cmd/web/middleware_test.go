package main

import (
	"bytes"
	"github.com/makarellav/codecapsule/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()
	defer rs.Body.Close()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "Content-Security-Policy",
			want: "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		},
		{
			name: "X-Content-Type-Options",
			want: "nosniff",
		},
		{
			name: "X-Frame-Options",
			want: "deny",
		},
		{
			name: "X-XSS-Protection",
			want: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, rs.Header.Get(tt.name), tt.want)
		})
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
