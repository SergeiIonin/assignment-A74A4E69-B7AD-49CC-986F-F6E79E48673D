package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain"
)

func TestTodoServiceImpl_GetTodos(t *testing.T) {
	todos := domain.Todos{Todos: []domain.Todo{
		{Id: 1, Todo: "Fix bug", Completed: false, UserId: 1},
		{Id: 2, Todo: "Write docs", Completed: true, UserId: 1},
	}}
	todosJSON, _ := json.Marshal(todos)

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		closeServer bool
		wantTodos   domain.Todos
		wantErr     bool
	}{
		{
			name: "valid JSON response decodes todos",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write(todosJSON)
			},
			wantTodos: todos,
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

			svc := NewTodoServiceImpl(url, client)
			got, err := svc.GetTodos(context.Background(), 1)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.wantTodos) {
				t.Errorf("got %+v, want %+v", got, tc.wantTodos)
			}
		})
	}
}
