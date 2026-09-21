package shopee

import (
	"context"
	"net/http"
)

type Token struct {
	AccessToken  string  `json:"access_token"`
	ExpireIn     int64   `json:"expire_in"`
	RefreshToken string  `json:"refresh_token"`
	RequestID    string  `json:"request_id,omitempty"`
	ShopID       int64   `json:"shop_id,omitempty"`
	ShopIDList   []int64 `json:"shop_id_list,omitempty"`
}

func (c *Client) GetAccessToken(ctx context.Context, code, shopID string) (Token, error) {
	var token Token
	err := c.Do(ctx, Request{
		Method: http.MethodPost,
		Path:   "/auth/token/get",
		ShopID: shopID,
		Body: map[string]any{
			"partner_id": c.cfg.PartnerID,
			"code":       code,
			"shop_id":    jsonNumber(shopID),
		},
	}, &token)
	return token, err
}

func (c *Client) RefreshAccessToken(ctx context.Context, refreshToken, shopID string) (Token, error) {
	var token Token
	err := c.Do(ctx, Request{
		Method: http.MethodPost,
		Path:   "/auth/access_token/get",
		ShopID: shopID,
		Body: map[string]any{
			"partner_id":    c.cfg.PartnerID,
			"refresh_token": refreshToken,
			"shop_id":       jsonNumber(shopID),
		},
	}, &token)
	return token, err
}

func jsonNumber(value string) any {
	var n int64
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return value
		}
		n = n*10 + int64(ch-'0')
	}
	return n
}
