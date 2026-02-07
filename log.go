package nod

import "log/slog"


type Client struct {
	log *slog.Logger
}


func New(logger *slog.Logger) *Client {
	if logger == nil {
		logger = slog.Default()
	}
	return &Client{
		log: logger,
	}
}

func SafePtrValue[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}
