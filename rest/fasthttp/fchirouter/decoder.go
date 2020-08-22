package fchirouter

import (
	"github.com/swaggest/fchi"
	"github.com/valyala/fasthttp"
	"net/url"
)

func PathToURLValues(rc *fasthttp.RequestCtx) (url.Values, error) {
	if routeCtx := fchi.RouteContext(rc); routeCtx != nil {
		params := make(url.Values, len(routeCtx.URLParams.Keys))

		for i, key := range routeCtx.URLParams.Keys {
			value, err := url.PathUnescape(routeCtx.URLParams.Values[i])
			if err != nil {
				return nil, err
			}

			params[key] = []string{value}
		}

		return params, nil
	}

	return nil, nil
}
