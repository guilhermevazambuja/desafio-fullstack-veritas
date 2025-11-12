package main

import "errors"

var tasks = []Task{
	{ID: "1", Title: "Clean Room", Completed: false},
	{ID: "2", Title: "Read Book", Completed: false},
	{ID: "3", Title: "Record Video", Completed: false},
}

func getTaskById(id string) (*Task, error) {
	for i, t := range tasks {
		if t.ID == id {
			return &tasks[i], nil
		}

	}
	return nil, ErrTaskNotFound
}
