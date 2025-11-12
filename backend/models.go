package main

type Task struct {
	ID        *string `json:"id"`
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

type SuccessResponse[T Task | []Task] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
