package main

import (
	"context"

	"github.com/fastly/compute-sdk-go/edgedict"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/auth"
	"github.com/leedo/activitypub-at-edge/proxy"
	"github.com/leedo/activitypub-at-edge/render"
)

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		d, err := edgedict.Open("auth")
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadGateway)
			w.Write([]byte(err.Error()))
		}

		clientId, err := d.Get("clientId")
		secret, err := d.Get("secret")

		a := auth.NewAuth(clientId, secret)

		switch r.URL.Path {
		case "/login":
			render.Login(w)
			return
		case "/logout":
			w.Header().Set("Set-Cookie", "auth=")
			w.Header().Set("Location", "/login")
			w.WriteHeader(fsthttp.StatusFound)
			return
		case "/gh_login":
			a.AuthHandler(ctx, w, r)
			return
		case "/oauth_callback":
			a.OAuthCallbackHandler(ctx, w, r)
			return
		}

		a.SetToken(r)

		if err := a.Check(ctx); err != nil {
			w.WriteHeader(fsthttp.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		p := proxy.NewProxy(w, a)

		switch r.URL.Path {
		case "/user":
			p.UserHandler(ctx, r)
		default:
			p.GenericRequestHandler(ctx, r)
		}
	})
}
