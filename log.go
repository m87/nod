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

func PtrVal[T any](p *T) any {
	if p == nil {
		return nil
	}
	return *p
}
