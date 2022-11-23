package main

import (
	"context"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/proxy"
)

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		proxy.NewProxy(w).GenericRequestHandler(ctx, r)
	})
}
