package main

import "errors"

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrInvalidPayload = errors.New("invalid task payload")
)
