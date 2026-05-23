package domain

type Todo struct {
	Id        int    `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"user_id"`
}

type Todos struct {
	Todos []Todo `json:"todos"`
}
