package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SergeiIonin/assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D/internal/domain"
	"github.com/gin-gonic/gin"
)

type DashboardService interface {
	GetDashboard(ctx context.Context, id int) (domain.Dashboard, error)
}

type Handlers struct {
	dashboardService DashboardService
}

func NewHandlers(dashboardService DashboardService) *Handlers {
	return &Handlers{
		dashboardService: dashboardService,
	}
}

func (h *Handlers) DashboardHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid ID format: %v", err)})
		return
	}
	fmt.Printf("id of request %d\n", id)
	dashboard, err := h.dashboardService.GetDashboard(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get dashboard: %v", err)})
		return
	}
	c.JSON(http.StatusOK, dashboard)
}

func (h *Handlers) ExecutionTimeHandler(c *gin.Context) {
	start := time.Now()

	c.Next()

	latency := time.Since(start)
	log.Printf("Request to %s took %v", c.Request.URL.Path, latency)
}

func getIdFromPath(path string) (int, error) {
	var id int
	_, err := fmt.Sscanf(path, "/dashboard/%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid path format: %v", err)
	}
	return id, nil
}
