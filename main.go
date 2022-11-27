package main

import (
	"context"

	"github.com/fastly/compute-sdk-go/edgedict"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/oauth"
	"github.com/leedo/activitypub-at-edge/render"
	"github.com/leedo/activitypub-at-edge/server"
)

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		d, err := edgedict.Open("oauth")
		if err != nil {
			w.WriteHeader(fsthttp.StatusBadGateway)
			w.Write([]byte(err.Error()))
		}

		clientId, err := d.Get("clientId")
		secret, err := d.Get("secret")

		a := oauth.NewOAuth(clientId, secret)

		switch r.URL.Path {
		case "/favicon.ico":
			w.WriteHeader(fsthttp.StatusNotFound)
			return
		case "/login":
			render.Login(w)
			return
		case "/logout":
			w.Header().Set("Set-Cookie", "auth=")
			w.Header().Set("Location", "/login")
			w.WriteHeader(fsthttp.StatusFound)
			return
		case "/gh_login":
			a.OAuthHandler(ctx, w, r)
			return
		case "/oauth_callback":
			a.OAuthCallbackHandler(ctx, w, r)
			return
		}

		a.SetToken(r)

		if err := a.Check(ctx); err != nil {
			//w.Header().Set("Location", "/login")
			w.WriteHeader(fsthttp.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		s := server.NewServer(w, a)

		switch r.URL.Path {
		case "/user":
			s.UserHandler(ctx, r)
		default:
			s.GenericRequestHandler(ctx, r)
		}
	})
}
