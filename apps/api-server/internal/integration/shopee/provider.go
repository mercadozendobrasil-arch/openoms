// Package shopee implements the Shopee marketplace integration.
package shopee

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openoms-org/openoms/apps/api-server/internal/integration"
	"github.com/openoms-org/openoms/apps/api-server/internal/model"
	shopeesdk "github.com/openoms-org/openoms/packages/shopee-go-sdk"
)

func init() {
	integration.RegisterMarketplaceProvider("shopee", func(credentials json.RawMessage, settings json.RawMessage) (integration.MarketplaceProvider, error) {
		return NewProvider(credentials, settings)
	})
}

type Credentials struct {
	PartnerID    int64  `json:"partner_id"`
	PartnerKey   string `json:"partner_key"`
	ShopID       string `json:"shop_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenExpiry  string `json:"token_expiry"`
	Sandbox      bool   `json:"sandbox,omitempty"`
}

type Provider struct {
	client *shopeesdk.Client
	creds  Credentials
}

func NewProvider(credentials json.RawMessage, _ json.RawMessage) (*Provider, error) {
	var creds Credentials
	if err := json.Unmarshal(credentials, &creds); err != nil {
		return nil, fmt.Errorf("shopee: parse credentials: %w", err)
	}
	if creds.PartnerID == 0 || creds.PartnerKey == "" || creds.ShopID == "" {
		return nil, fmt.Errorf("shopee: partner_id, partner_key and shop_id are required")
	}
	baseURL := shopeesdk.ProductionBaseURL
	if creds.Sandbox {
		baseURL = shopeesdk.SandboxBaseURL
	}
	return &Provider{
		client: shopeesdk.NewClient(shopeesdk.Config{BaseURL: baseURL, PartnerID: creds.PartnerID, PartnerKey: creds.PartnerKey}),
		creds: creds,
	}, nil
}

func (p *Provider) ProviderName() string { return "shopee" }
func (p *Provider) SDKClient() *shopeesdk.Client { return p.client }\nfunc (p *Provider) AccessToken() string { return p.creds.AccessToken }\nfunc (p *Provider) ShopID() string { return p.creds.ShopID }

func (p *Provider) ShopInfo(ctx context.Context) (shopeesdk.ShopInfo, error) {
	return p.client.GetShopInfo(ctx, p.creds.AccessToken, p.creds.ShopID)
}

// Marketplace methods are enabled incrementally as each Shopee module is ported.
func (p *Provider) PollOrders(context.Context, string) ([]integration.MarketplaceOrder, string, error) {
	return nil, "", fmt.Errorf("shopee: order sync not implemented yet")
}
func (p *Provider) GetOrder(context.Context, string) (*integration.MarketplaceOrder, error) {
	return nil, fmt.Errorf("shopee: order sync not implemented yet")
}
func (p *Provider) PushOffer(context.Context, *model.Product, map[string]any) (string, error) {
	return "", fmt.Errorf("shopee: product publishing not implemented yet")
}
func (p *Provider) UpdateStock(context.Context, string, int) error {
	return fmt.Errorf("shopee: inventory sync not implemented yet")
}
func (p *Provider) UpdatePrice(context.Context, string, float64) error {
	return fmt.Errorf("shopee: price sync not implemented yet")
}

func (p *Provider) TokenExpiry() (time.Time, error) {
	if p.creds.TokenExpiry == "" { return time.Time{}, nil }
	return time.Parse(time.RFC3339, p.creds.TokenExpiry)
}
