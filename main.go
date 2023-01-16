package main

import (
	"context"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/config"
	"github.com/leedo/activitypub-at-edge/oauth"
	"github.com/leedo/activitypub-at-edge/render"
	"github.com/leedo/activitypub-at-edge/server"
	"github.com/leedo/activitypub-at-edge/store"
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
		c, err := config.ReadConfig()
		if err != nil {
			server.ErrorPage(fsthttp.StatusInternalServerError, w, "config error: "+err.Error())
			return
		}

		a := oauth.NewOAuth(c.OAuthClientId, c.OAuthSecret)

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

		t, err := store.Open(c.StoreName)
		if err != nil {
			server.ErrorPage(fsthttp.StatusInternalServerError, w, "store error: "+err.Error())
			return
		}

		u.SetStore(t)
		s := server.NewServer(u, t)

		if err := u.PreloadSettings(); err != nil {
			server.ErrorPage(fsthttp.StatusInternalServerError, w, "settings error: "+err.Error())
			return
		}

		if r.Method == "GET" {
			switch r.URL.Path {
			case logoutPath:
				a.OAuthLogoutHandler(ctx, w, r, u, loginPath)
			case "/user", "/":
				s.UserHandler(ctx, w, r)
			default:
				s.GenericRequestHandler(ctx, w, r)
			}
		}
		if r.Method == "POST" && r.URL.Path == "/subscriptions" {
			s.SubscribeHandler(ctx, w, r)
			return
		}
	})
}
