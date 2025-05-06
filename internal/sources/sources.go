package sources

import "context"

type Redirect struct {
	URLPattern string `json:"urlPattern"`
}

type Source interface {
	GetRedirectForKey(ctx context.Context, key string) (Redirect, error)
	GetAllRedirects(ctx context.Context) (map[string]Redirect, error)
}
