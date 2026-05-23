package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain"
	"github.com/gin-gonic/gin"
)

type mockDashboardService struct {
	dashboard domain.Dashboard
	err       error
}

func (m *mockDashboardService) GetDashboard(_ context.Context, _ int) (domain.Dashboard, error) {
	return m.dashboard, m.err
}

func TestDashboardHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dashboard := domain.Dashboard{
		Id:               1,
		FullName:         "John Doe",
		Status:           "Veteran",
		PendingTaskCount: 2,
		NextUrgentTask:   "Fix bug",
	}

	tests := []struct {
		name       string
		id         string
		svc        DashboardService
		wantStatus int
		wantBody   *domain.Dashboard
	}{
		{
			name:       "non-numeric id returns 400",
			id:         "abc",
			svc:        &mockDashboardService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error returns 500",
			id:         "1",
			svc:        &mockDashboardService{err: errors.New("internal error")},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "valid request returns 200 with dashboard JSON",
			id:         "1",
			svc:        &mockDashboardService{dashboard: dashboard},
			wantStatus: http.StatusOK,
			wantBody:   &dashboard,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			h := NewHandlers(tc.svc)
			router.GET("/dashboard/:id", h.DashboardHandler)

			req := httptest.NewRequest(http.MethodGet, "/dashboard/"+tc.id, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantBody != nil {
				var got domain.Dashboard
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode response body: %v", err)
				}
				if got != *tc.wantBody {
					t.Errorf("body = %+v, want %+v", got, *tc.wantBody)
				}
			}
		})
	}
}

func TestGetIdFromPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantId  int
		wantErr bool
	}{
		{
			name:   "valid path returns id",
			path:   "/dashboard/42",
			wantId: 42,
		},
		{
			name:    "non-numeric segment returns error",
			path:    "/dashboard/abc",
			wantErr: true,
		},
		{
			name:    "empty string returns error",
			path:    "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getIdFromPath(tc.path)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantId {
				t.Errorf("got %d, want %d", got, tc.wantId)
			}
		})
	}
}
