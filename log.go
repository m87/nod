package nod

import "log/slog"

// Client provides a configured logger for nod operations.
type Client struct {
	log *slog.Logger
}

// New creates a new Client with the given logger. If logger is nil, the default logger is used.
func New(logger *slog.Logger) *Client {
	if logger == nil {
		logger = slog.Default()
	}
	return &Client{
		log: logger,
	}
}

// SafePtrValue dereferences a pointer, returning the zero value of T if ptr is nil.
func SafePtrValue[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}
