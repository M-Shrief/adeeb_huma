package cache

import (
	"adeeb_huma/config"

	glide "github.com/valkey-io/valkey-glide/go/v2"
	glide_config "github.com/valkey-io/valkey-glide/go/v2/config"
)

var Client *glide.Client

// Create a new Cache ValKey Client
func NewCache() (*glide.Client, error) {
	cache_config := glide_config.NewClientConfiguration().
		WithAddress(&glide_config.NodeAddress{Host: config.Cache_HOST, Port: config.Cache_PORT}).
		WithCredentials(glide_config.NewServerCredentials("", config.Cache_PASSWORD))
	var err error
	Client, err = glide.NewClient(cache_config)
	if err != nil {
		return nil, err
	}

	return Client, nil
}
