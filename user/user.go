package user

import (
	"fmt"
	"strings"

	"github.com/leedo/activitypub-at-edge/store"
	"github.com/valyala/fastjson"
)

const (
	settingsKey = "settings"
)

type User struct {
	OauthToken string
	Login      string
	Id         string
	AvatarUrl  string
	settings   *Settings
	s          *store.Store
}

type Settings struct {
	Subscriptions []string
}

func (u *User) Lock() error {
	err := u.s.Lock(u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Unlock() error {
	if err := u.s.Unlock(u.Id); err != nil {
		return err
	}
	return nil
}

func (u *User) SetStore(s *store.Store) {
	u.s = s
}

func (u *User) PreloadSettings() error {
	s, err := u.readSettings()
	if err != nil {
		if err == store.KeyNotFound {
			u.settings = &Settings{}
			return nil
		}
		return err
	}
	u.settings = s
	return nil
}

func (u *User) Settings() (*Settings, error) {
	if u.settings != nil {
		return u.settings, nil
	}

	s, err := u.readSettings()
	if err != nil {
		return nil, err
	}

	u.settings = s
	return u.settings, nil
}

func (u *User) saveSettings() error {
	if u.settings == nil {
		return nil
	}

	if err := u.s.Insert(u.storeKey(settingsKey), u.settings.String()); err != nil {
		return err
	}

	return nil
}

func (s *Settings) String() string {
	var out string

	out += `{`
	out += `"subscriptions":`
	out += `[`

	var subs []string
	for _, f := range s.Subscriptions {
		subs = append(subs, `"`+f+`"`)
	}

	out += strings.Join(subs, ",")
	out += `]`
	out += `}`

	return out
}

func (s *Settings) IsSubscribed(url string) bool {
	for _, v := range s.Subscriptions {
		if v == url {
			return true
		}
	}

	return false
}
func (u *User) Subscribe(url string) error {
	settings, err := u.Settings()
	if err != nil {
		return err
	}

	if settings.IsSubscribed(url) {
		return fmt.Errorf("already subscribed %s", url)
	}

	settings.Subscriptions = append(settings.Subscriptions, url)
	return u.saveSettings()
}

func (u *User) Unsubscribe(url string) error {
	settings, err := u.Settings()
	if err != nil {
		return err
	}

	var newSubs []string
	for _, f := range settings.Subscriptions {
		if f != url {
			newSubs = append(newSubs, f)
		}
	}
	settings.Subscriptions = newSubs
	return u.saveSettings()
}

func (u *User) readSettings() (*Settings, error) {
	s, err := u.s.Lookup(u.storeKey(settingsKey))
	if err != nil {
		return nil, err
	}

	var p fastjson.Parser
	v, err := p.ParseBytes([]byte(s))
	if err != nil {
		return nil, err
	}

	settings := &Settings{}
	for _, f := range v.GetArray("subscriptions") {
		settings.Subscriptions = append(settings.Subscriptions, string(f.GetStringBytes()))
	}

	return settings, nil
}

func (u *User) storeKey(key string) string {
	return u.Id + "/" + key
}
