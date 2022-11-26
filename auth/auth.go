package auth

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/valyala/fastjson"
)

const host = "https://github.com"
const apiHost = "https://api.github.com"

type Auth struct {
	clientId string
	secret   string
	token    string
	User     *User
}

type User struct {
	Login  string
	Id     string
	Avatar string
}

func NewAuth(clientId, secret string) *Auth {
	return &Auth{
		clientId: clientId,
		secret:   secret,
	}
}

func (a *Auth) SetToken(r *fsthttp.Request) {
	c, err := r.Cookie("auth")
	if err != nil {
		return
	}
	a.token = c.Value
}

func (a *Auth) Check(ctx context.Context) error {
	buf := []byte(`{"access_token":"` + a.token + `"}`)
	b := bytes.NewBuffer(buf)
	req, err := fsthttp.NewRequest("POST", apiHost+"/applications/"+a.clientId+"/token", b)
	if err != nil {
		return err
	}

	s := base64.StdEncoding.EncodeToString([]byte(a.clientId + ":" + a.secret))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+s)
	req.Header.Set("User-Agent", "activitypub-at-edge")

	resp, err := req.Send(ctx, req.URL.Host)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != fsthttp.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	j, err := fastjson.Parse(string(body))
	if err != nil {
		return err
	}

	u := j.Get("user")

	a.User = &User{
		Login:  string(u.GetStringBytes("login")),
		Id:     string(u.GetStringBytes("id")),
		Avatar: string(u.GetStringBytes("avatar")),
	}

	return nil
}

func (a *Auth) AuthHandler(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
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

func (a *Auth) Error(w fsthttp.ResponseWriter, msg string) {
	w.WriteHeader(fsthttp.StatusForbidden)
	w.Write([]byte(msg))
}

func (a *Auth) OAuthCallbackHandler(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		a.Error(w, "oauth failure")
		return
	}

	buf := []byte(fmt.Sprintf(`{"client_id":"%s","code":"%s","client_secret":"%s"}`, a.clientId, code, a.secret))
	b := bytes.NewBuffer(buf)
	req, err := fsthttp.NewRequest("POST", host+"/login/oauth/access_token", b)
	if err != nil {
		a.Error(w, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := req.Send(ctx, req.URL.Host)
	if err != nil {
		a.Error(w, err.Error())
		return
	}

	defer resp.Body.Close()

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		a.Error(w, err.Error())
		return
	}

	j, err := fastjson.Parse(string(d))
	if err != nil {
		a.Error(w, err.Error())
		return
	}

	t := string(j.GetStringBytes("access_token"))
	if t == "" {
		a.Error(w, "no access token")
		return
	}

	w.Header().Set("Set-Cookie", "auth="+t)
	w.Header().Set("Location", "/")
	w.WriteHeader(fsthttp.StatusFound)
}
