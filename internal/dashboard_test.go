package internal

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain"
)

type mockProfileService struct {
	profile domain.Profile
	err     error
}

func (m *mockProfileService) GetProfile(_ context.Context, _ int) (domain.Profile, error) {
	return m.profile, m.err
}

type mockTodosService struct {
	todos domain.Todos
	err   error
}

func (m *mockTodosService) GetTodos(_ context.Context, _ int) (domain.Todos, error) {
	return m.todos, m.err
}

func TestDashboardImpl_GetDashboard(t *testing.T) {
	veteran := domain.Profile{Id: 1, FirstName: "John", LastName: "Doe", Age: 55}
	rookie := domain.Profile{Id: 2, FirstName: "Jane", LastName: "Smith", Age: 30}

	mixedTodos := domain.Todos{Todos: []domain.Todo{
		{Id: 1, Todo: "Fix bug", Completed: false, UserId: 1},
		{Id: 2, Todo: "Write docs", Completed: true, UserId: 1},
		{Id: 3, Todo: "Add tests", Completed: false, UserId: 1},
	}}
	allCompletedTodos := domain.Todos{Todos: []domain.Todo{
		{Id: 1, Todo: "Done task", Completed: true, UserId: 2},
	}}

	todosWarning := "Todos Unavailable"

	tests := []struct {
		name          string
		profileSvc    ProfileService
		todosSvc      TodosService
		id            int
		wantDashboard domain.Dashboard
		wantErr       bool
	}{
		{
			name:       "profile service error propagates as error",
			profileSvc: &mockProfileService{err: errors.New("profile unavailable")},
			todosSvc:   &mockTodosService{todos: mixedTodos},
			id:         1,
			wantErr:    true,
		},
		{
			name:       "todos service error returns partial dashboard with warning",
			profileSvc: &mockProfileService{profile: veteran},
			todosSvc:   &mockTodosService{err: errors.New("todos unavailable")},
			id:         1,
			wantDashboard: domain.Dashboard{
				Id:           1,
				FullName:     "John Doe",
				Status:       "Veteran",
				ErrorWarning: &todosWarning,
			},
		},
		{
			name:       "both succeed with mixed todos returns full dashboard",
			profileSvc: &mockProfileService{profile: veteran},
			todosSvc:   &mockTodosService{todos: mixedTodos},
			id:         1,
			wantDashboard: domain.Dashboard{
				Id:               1,
				FullName:         "John Doe",
				Status:           "Veteran",
				PendingTaskCount: 2,
				NextUrgentTask:   "Fix bug",
			},
		},
		{
			name:       "age under 50 sets Rookie status",
			profileSvc: &mockProfileService{profile: rookie},
			todosSvc:   &mockTodosService{todos: allCompletedTodos},
			id:         2,
			wantDashboard: domain.Dashboard{
				Id:               2,
				FullName:         "Jane Smith",
				Status:           "Rookie",
				PendingTaskCount: 0,
				NextUrgentTask:   "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewDashboardImpl(tc.profileSvc, tc.todosSvc)
			got, err := svc.GetDashboard(context.Background(), tc.id)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.wantDashboard) {
				t.Errorf("got %+v, want %+v", got, tc.wantDashboard)
			}
		})
	}
}
