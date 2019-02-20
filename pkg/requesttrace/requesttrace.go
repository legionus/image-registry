package requesttrace

import (
	"context"
	"fmt"
	"net/http"

	dcontext "github.com/docker/distribution/context"
)

const (
	requestHeader = "X-Registry-Request-URL"
)

type requestTracer struct {
	ctx context.Context
	req *http.Request
}

func New(ctx context.Context, req *http.Request) *requestTracer {
	return &requestTracer{
		ctx: ctx,
		req: req,
	}
}

func (rt *requestTracer) ModifyRequest(req *http.Request) (err error) {
	for _, k := range rt.req.Header[requestHeader] {
		if k == req.URL.String() {
			err = fmt.Errorf("Request to %q is denied because a loop is detected", k)
			dcontext.GetLogger(rt.ctx).Error(err.Error())
			return
		}
		req.Header.Add(requestHeader, k)
	}
	req.Header.Add(requestHeader, req.URL.String())
	return nil
}
