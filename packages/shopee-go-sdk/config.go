package shopee

import (
	"net/http"
	"time"
)

const (
	ProductionBaseURL = "https://partner.shopeemobile.com"
	SandboxBaseURL    = "https://openplatform.sandbox.test-stable.shopee.sg"
)

type Config struct {
	BaseURL    string
	APIVersion string
	PartnerID  int64
	PartnerKey string
	HTTPClient *http.Client
	Now        func() time.Time
}

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = ProductionBaseURL
	}
	if c.APIVersion == "" {
		c.APIVersion = "v2"
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	if c.Now == nil {
		c.Now = time.Now
	}
	return c
}
