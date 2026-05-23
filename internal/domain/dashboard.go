package domain

type Dashboard struct {
	Id               int     `json:"id"`
	FullName         string  `json:"full_name"`
	Status           string  `json:"status"`
	PendingTaskCount int     `json:"pending_task_count"`
	NextUrgentTask   string  `json:"next_urgent_task"`
	ErrorWarning     *string `json:"error_warning"`
}
