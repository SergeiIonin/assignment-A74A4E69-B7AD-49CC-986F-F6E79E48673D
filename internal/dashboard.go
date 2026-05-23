package internal

import (
	"context"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain"
	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain/fetcher"
)

type ProfileService interface {
	GetProfile(ctx context.Context, id int) (domain.Profile, error)
}

type TodosService interface {
	GetTodos(ctx context.Context, id int) (domain.Todos, error)
}

type DashboardImpl struct {
	profileService ProfileService
	todosService   TodosService
}

func NewDashboardImpl(profileService ProfileService, todosService TodosService) *DashboardImpl {
	return &DashboardImpl{
		profileService: profileService,
		todosService:   todosService,
	}
}

func (d *DashboardImpl) GetDashboard(ctx context.Context, id int) (domain.Dashboard, error) {
	profileResChan := fetcher.Fetch(ctx, func(ctx context.Context) fetcher.Result[domain.Profile] {
		profile, err := d.profileService.GetProfile(ctx, id)
		return fetcher.Result[domain.Profile]{Data: profile, Error: err}
	})

	todosResChan := fetcher.Fetch(ctx, func(ctx context.Context) fetcher.Result[domain.Todos] {
		todos, err := d.todosService.GetTodos(ctx, id)
		return fetcher.Result[domain.Todos]{Data: todos, Error: err}
	})

	profileResult := <-profileResChan
	todosResult := <-todosResChan

	profile := profileResult.Data
	if profileResult.Error != nil {
		return domain.Dashboard{}, profileResult.Error
	}

	todos := todosResult.Data
	if todosResult.Error != nil {
		return d.composePartialDashboard(profile), nil
	}

	return d.composeDashboard(profile, todos.Todos), nil
}

func (d *DashboardImpl) composePartialDashboard(profile domain.Profile) domain.Dashboard {
	fullName := getFullName(profile)
	status := getStatus(profile)
	id := profile.Id

	warning := "Todos Unavailable"

	return domain.Dashboard{
		Id:           id,
		FullName:     fullName,
		Status:       status,
		ErrorWarning: &warning,
	}
}

func (d *DashboardImpl) composeDashboard(profile domain.Profile, todos []domain.Todo) domain.Dashboard {
	fullName := getFullName(profile)
	status := getStatus(profile)
	id := profile.Id
	pendingTaskCount := 0
	var nextUrgentTask string
	for _, todo := range todos {
		if !todo.Completed {
			pendingTaskCount++
			if nextUrgentTask == "" {
				nextUrgentTask = todo.Todo
			}
		}
	}

	return domain.Dashboard{
		Id:               id,
		FullName:         fullName,
		Status:           status,
		PendingTaskCount: pendingTaskCount,
		NextUrgentTask:   nextUrgentTask,
	}
}

func getFullName(profile domain.Profile) string {
	return profile.FirstName + " " + profile.LastName
}

func getStatus(profile domain.Profile) string {
	if profile.Age < 50 {
		return "Rookie"
	}
	return "Veteran"
}
