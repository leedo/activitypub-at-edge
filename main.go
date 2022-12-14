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
		case oauthPath:
			a.OAuthHandler(ctx, w, r)
			return
		case oauthCallbackPath:
			a.OAuthCallbackHandler(ctx, w, r)
			return
		}

		u, err := a.Check(ctx, r)
		if err != nil {
			w.Header().Set("Location", loginPath)
			w.WriteHeader(fsthttp.StatusFound)
			return
		}

		s := server.NewServer(u)

		switch r.URL.Path {
		case logoutPath:
			a.OAuthLogoutHandler(ctx, w, r, u, loginPath)
		case "/user":
			s.UserHandler(ctx, w, r)
		default:
			s.GenericRequestHandler(ctx, w, r)
		}
	})
}
