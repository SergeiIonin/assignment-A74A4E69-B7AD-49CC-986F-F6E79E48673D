package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain"
)

type TodoServiceImpl struct {
	TodosPath  string
	httpClient *http.Client
}

func NewTodoServiceImpl(todosPath string, httpClient *http.Client) *TodoServiceImpl {
	return &TodoServiceImpl{
		TodosPath:  todosPath,
		httpClient: httpClient,
	}
}

func (t *TodoServiceImpl) GetTodos(ctx context.Context, id int) (domain.Todos, error) {
	log.Printf("Fetching todos for path: %s", t.TodosPath+"/"+strconv.Itoa(id))
	resp, err := t.httpClient.Get(t.TodosPath + "/" + strconv.Itoa(id))
	if err != nil {
		log.Printf("Error fetching todos: %v", err)
		return domain.Todos{}, err
	}
	defer resp.Body.Close()

	var todos domain.Todos
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return domain.Todos{}, err
	}

	return todos, nil
}
