package shopee

import (
	"context"
	"net/http"
)

type ShopInfo struct {
	ShopID      int64  `json:"shop_id"`
	ShopName    string `json:"shop_name"`
	Region      string `json:"region,omitempty"`
	Status      string `json:"status,omitempty"`
	SIPAffiliate bool  `json:"sip_affiliate,omitempty"`
}

type ShopProfile struct {
	ShopLogo        string `json:"shop_logo,omitempty"`
	Description     string `json:"description,omitempty"`
	ShopName        string `json:"shop_name,omitempty"`
}

func (c *Client) GetShopInfo(ctx context.Context, accessToken, shopID string) (ShopInfo, error) {
	var out ShopInfo
	err := c.Do(ctx, Request{Method: http.MethodGet, Path: "/shop/get_shop_info", AccessToken: accessToken, ShopID: shopID}, &out)
	return out, err
}

func (c *Client) GetProfile(ctx context.Context, accessToken, shopID string) (ShopProfile, error) {
	var out ShopProfile
	err := c.Do(ctx, Request{Method: http.MethodGet, Path: "/shop/get_profile", AccessToken: accessToken, ShopID: shopID}, &out)
	return out, err
}
