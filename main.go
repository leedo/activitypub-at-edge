package main

import (
	"context"

	"github.com/fastly/compute-sdk-go/edgedict"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/oauth"
	"github.com/leedo/activitypub-at-edge/render"
	"github.com/leedo/activitypub-at-edge/server"
)

const (
	loginPath         = "/login"
	logoutPath        = "/logout"
	oauthPath         = "/oauth_login"
	oauthCallbackPath = "/oauth_callback"
	faviconPath       = "/favicon.ico"
)

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		d, err := edgedict.Open("oauth")
		if err != nil {
			w.WriteHeader(fsthttp.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		clientId, err := d.Get("clientId")
		secret, err := d.Get("secret")

		a := oauth.NewOAuth(clientId, secret)

		switch r.URL.Path {
		case faviconPath:
			w.WriteHeader(fsthttp.StatusNotFound)
			return
		case loginPath:
			render.Login(w)
			return
		case logoutPath:
			w.Header().Set("Set-Cookie", "auth=")
			w.Header().Set("Location", loginPath)
			w.WriteHeader(fsthttp.StatusFound)
			return
		case oauthPath:
			a.OAuthHandler(ctx, w, r)
			return
		case oauthCallbackPath:
			a.OAuthCallbackHandler(ctx, w, r)
			return
		}

		a.SetToken(r)

		if err := a.Check(ctx); err != nil {
			w.Header().Set("Location", "/login")
			w.WriteHeader(fsthttp.StatusFound)
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
