package main

type Task struct {
	ID     *string `json:"id"`
	Title  *string `json:"title"`
	Status *string `json:"status"` // "to_do", "in_progress", "done"
}

type SuccessResponse[T Task | []Task] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
