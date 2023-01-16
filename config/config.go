package config

import "github.com/fastly/compute-sdk-go/edgedict"

type Config struct {
	OAuthClientId string
	OAuthSecret   string
	StoreName     string
}

func ReadConfig() (*Config, error) {
	d, err := edgedict.Open("oauth")
	if err != nil {
		return nil, err
	}

	clientId, err := d.Get("clientId")
	if err != nil {
		return nil, err
	}
	secret, err := d.Get("secret")
	if err != nil {
		return nil, err
	}
	storeName, err := d.Get("storeName")
	if err != nil {
		return nil, err
	}
	return &Config{
		OAuthClientId: clientId,
		OAuthSecret:   secret,
		StoreName:     storeName,
	}, nil
}
