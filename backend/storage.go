package main

var tasks = []Task{
	{ID: strPtr("1"), Title: strPtr("Clean Room"), Status: strPtr("to_do")},
	{ID: strPtr("2"), Title: strPtr("Read Book"), Status: strPtr("in_progress")},
	{ID: strPtr("3"), Title: strPtr("Record Video"), Status: strPtr("done")},
}

func getTaskById(id string) (*Task, error) {
	for i, t := range tasks {
		if t.ID != nil && *t.ID == id {
			return &tasks[i], nil
		}

	}
	return nil, ErrTaskNotFound
}

func strPtr(s string) *string { return &s }
