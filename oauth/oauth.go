package oauth

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/leedo/activitypub-at-edge/user"
	"github.com/valyala/fastjson"
)

const (
	host     = "https://github.com"
	apiHost  = "https://api.github.com"
	jsonType = "application/json"
)

type OAuth struct {
	clientId string
	secret   string
}

func NewOAuth(clientId, secret string) *OAuth {
	return &OAuth{
		clientId: clientId,
		secret:   secret,
	}
}

func (a *OAuth) Check(ctx context.Context, r *fsthttp.Request) (*user.User, error) {
	t, err := r.Cookie("auth")
	if err != nil {
		return nil, err
	}

	body, err := a.checkToken(ctx, t.Value)
	if err != nil {
		return nil, err
	}

	j, err := fastjson.Parse(string(body))
	if err != nil {
		return nil, err
	}

	u := j.Get("user")

	return &user.User{
		OauthToken: t.Value,
		Login:      string(u.GetStringBytes("login")),
		Id:         fmt.Sprintf("%d", u.GetInt("id")),
		AvatarUrl:  string(u.GetStringBytes("avatar_url")),
	}, nil
}

func (a *OAuth) Error(w fsthttp.ResponseWriter, msg string) {
	w.WriteHeader(fsthttp.StatusForbidden)
	w.Write([]byte(msg))
}

func (a *OAuth) basicAuth() string {
	return base64.StdEncoding.EncodeToString([]byte(a.clientId + ":" + a.secret))
}

func (a *OAuth) request(ctx context.Context, method string, url string, b io.Reader) ([]byte, error) {
	req, err := fsthttp.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", jsonType)
	req.Header.Set("Accept", jsonType)
	req.Header.Set("Authorization", "Basic "+a.basicAuth())
	req.Header.Set("User-Agent", "activitypub-at-edge")

	resp, err := req.Send(ctx, req.URL.Host)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != fsthttp.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	return body, nil
}

func (a *OAuth) checkToken(ctx context.Context, token string) ([]byte, error) {
	buf := []byte(`{"access_token":"` + token + `"}`)
	b := bytes.NewBuffer(buf)
	return a.request(ctx, "POST", apiHost+"/applications/"+a.clientId+"/token", b)
}

func (a *OAuth) deleteToken(ctx context.Context, t string) ([]byte, error) {
	buf := []byte(fmt.Sprintf(`{"access_token":"%s"}`, t))
	b := bytes.NewBuffer(buf)
	return a.request(ctx, "DELETE", apiHost+"/applications/"+a.clientId+"/token", b)
}

func (a *OAuth) createToken(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("empty code")
	}

	buf := []byte(fmt.Sprintf(`{"client_id":"%s","code":"%s","client_secret":"%s"}`, a.clientId, code, a.secret))
	b := bytes.NewBuffer(buf)

	body, err := a.request(ctx, "POST", host+"/login/oauth/access_token", b)
	if err != nil {
		return "", err
	}

	j, err := fastjson.Parse(string(body))
	if err != nil {
		return "", err
	}

	t := string(j.GetStringBytes("access_token"))
	if t == "" {
		return "", fmt.Errorf("no token")
	}

	return t, nil
}

func (a *OAuth) OAuthHandler(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
	u, err := url.Parse(host + "/login/oauth/authorize")
	if err != nil {
		a.Error(w, err.Error())
	}

	v := url.Values{}
	v.Set("client_id", a.clientId)
	u.RawQuery = v.Encode()

	w.Header().Add("Location", u.String())
	w.WriteHeader(fsthttp.StatusFound)
}

func (a *OAuth) OAuthLogoutHandler(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request, u *user.User, redirect string) {
	if _, err := a.deleteToken(ctx, u.OauthToken); err != nil {
		a.Error(w, "oauth delete failure")
		return
	}

	w.Header().Set("Set-Cookie", "auth=")
	w.Header().Set("Location", redirect)
	w.WriteHeader(fsthttp.StatusFound)
}

func (a *OAuth) OAuthCallbackHandler(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
	t, err := a.createToken(ctx, r.URL.Query().Get("code"))
	if err != nil {
		a.Error(w, err.Error())
		return
	}

	w.Header().Set("Set-Cookie", "auth="+t)
	w.Header().Set("Location", "/")
	w.WriteHeader(fsthttp.StatusFound)
}
