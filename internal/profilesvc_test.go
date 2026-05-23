package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain"
)

func TestProfileServiceImpl_GetProfile(t *testing.T) {
	profile := domain.Profile{Id: 1, FirstName: "John", LastName: "Doe", Age: 55}
	profileJSON, _ := json.Marshal(profile)

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		closeServer bool
		wantProfile domain.Profile
		wantErr     bool
	}{
		{
			name: "valid JSON response decodes profile",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(profileJSON)
			},
			wantProfile: profile,
		},
		{
			name: "invalid JSON body returns decode error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("not-json"))
			},
			wantErr: true,
		},
		{
			name:        "connection refused returns transport error",
			handler:     func(w http.ResponseWriter, r *http.Request) {},
			closeServer: true,
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewServer(tc.handler)
			url := s.URL
			client := s.Client()

			if tc.closeServer {
				s.Close()
				client = http.DefaultClient
			} else {
				t.Cleanup(s.Close)
			}

			svc := NewProfileServiceImpl(url, client)
			got, err := svc.GetProfile(context.Background(), 1)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantProfile {
				t.Errorf("got %+v, want %+v", got, tc.wantProfile)
			}
		})
	}
}
