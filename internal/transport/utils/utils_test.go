package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseJSON(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name        string
		body        []byte
		wantErr     bool
		expected    payload
		expectedErr string
	}{
		{
			name:     "valid JSON",
			body:     []byte(`{"name":"Misha", "age":30}`),
			wantErr:  false,
			expected: payload{Name: "Misha", Age: 30},
		},
		{
			name:        "empty body",
			body:        nil,
			wantErr:     true,
			expectedErr: "The body is empty",
		},
		{
			name:        "malformed JSON",
			body:        []byte(`{"name":"Petr", "age":}`),
			wantErr:     true,
			expectedErr: "invalid character '}' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.body))

			if tt.body == nil {
				req.Body = nil
			}

			var got payload
			err := ParseJSON(req, &got)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if strings.Compare(err.Error(), tt.expectedErr) != 0 {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	type response struct {
		Message string
		Code    int
	}

	tests := []struct {
		name          string
		status        int
		payload       any
		wantBody      string
		wantHeaderKey string
		wantHeaderVal string
	}{
		{
			name:          "simple payload",
			status:        http.StatusOK,
			payload:       response{Message: "ok", Code: 0},
			wantBody:      `{"Message":"ok","Code":0}` + "\n",
			wantHeaderKey: "Content-Type",
			wantHeaderVal: "application/json",
		},
		{
			name:          "nil payload",
			status:        http.StatusNoContent,
			payload:       nil,
			wantBody:      "null\n",
			wantHeaderKey: "Content-Type",
			wantHeaderVal: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := WriteJSON(rr, tt.status, tt.payload)

			if err != nil {
				t.Fatalf("Unexpected error err:%v", err)
			}

			if rr.Code != tt.status {
				t.Fatalf("status = %d, expected %d", rr.Code, tt.status)
			}

			if got := rr.Header().Get(tt.wantHeaderKey); got != tt.wantHeaderVal {
				t.Fatalf("header %s = %q, expected %q", tt.wantHeaderKey, got, tt.wantHeaderVal)
			}

			if rr.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rr.Body.String(), tt.wantBody)
			}
		})
	}
}
