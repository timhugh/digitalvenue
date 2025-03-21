package db

import "context"

type Client interface {
	ExecuteQuery(ctx context.Context, query string, args ...any) Result
}
