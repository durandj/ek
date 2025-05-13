package sources

import (
	"context"
)

type InMemorySource struct {
	redirects map[string]Redirect
}

var _ Source = (*InMemorySource)(nil)

type InMemorySourceParams struct {
	Redirects map[string]Redirect
}

func NewInMemorySource(params InMemorySourceParams) *InMemorySource {
	return &InMemorySource{
		redirects: params.Redirects,
	}
}

func (src *InMemorySource) GetRedirectForKey(ctx context.Context, key string) (Redirect, error) {
	redirect, ok := src.redirects[key]
	if !ok {
		return Redirect{}, errNoSuchRedirectKey
	}

	return redirect, nil
}

func (src *InMemorySource) GetAllRedirects(ctx context.Context) (map[string]Redirect, error) {
	return src.redirects, nil
}
