package main

type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ListResp struct {
	Data []Task `json:"data"`
}
