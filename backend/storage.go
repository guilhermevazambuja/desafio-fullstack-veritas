package main

var tasks = []Task{
	{ID: strPtr("1"), Title: strPtr("Clean Room"), Completed: boolPtr(false)},
	{ID: strPtr("2"), Title: strPtr("Read Book"), Completed: boolPtr(false)},
	{ID: strPtr("3"), Title: strPtr("Record Video"), Completed: boolPtr(false)},
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
func boolPtr(b bool) *bool    { return &b }
