package api

import (
	"fmt"
	"jcloud/config"
)

type API struct {
	cfg *config.Config
}

func New(cfg *config.Config) *API {
	return &API{cfg: cfg}
}

func (a *API) MakeRequest() {
	fmt.Printf("Making request with API key: %s\n", a.cfg.APIKey)
}
