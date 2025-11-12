package main

import "errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidPayload    = errors.New("invalid task payload")
	ErrIncompletePayload = errors.New("all fields must be provided")
	ErrIDMismatch        = errors.New("payload ID does not match URL ID")
)
