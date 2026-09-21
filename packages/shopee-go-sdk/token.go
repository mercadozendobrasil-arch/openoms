package shopee

import (
	"context"
	"errors"
	"time"
)

var ErrTokenUnavailable = errors.New("shopee: token unavailable")

type StoredToken struct {
	ShopID       string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type TokenStore interface {
	Load(ctx context.Context, shopID string) (StoredToken, error)
	Save(ctx context.Context, token StoredToken) error
}

type TokenManager struct {
	client *Client
	store  TokenStore
	leeway time.Duration
	now    func() time.Time
}

func NewTokenManager(client *Client, store TokenStore, leeway time.Duration) *TokenManager {
	if leeway <= 0 {
		leeway = 5 * time.Minute
	}
	return &TokenManager{client: client, store: store, leeway: leeway, now: client.cfg.Now}
}

func (m *TokenManager) AccessToken(ctx context.Context, shopID string) (string, error) {
	stored, err := m.store.Load(ctx, shopID)
	if err != nil {
		return "", err
	}
	if stored.AccessToken == "" || stored.RefreshToken == "" {
		return "", ErrTokenUnavailable
	}
	if m.now().Add(m.leeway).Before(stored.ExpiresAt) {
		return stored.AccessToken, nil
	}

	refreshed, err := m.client.RefreshAccessToken(ctx, stored.RefreshToken, shopID)
	if err != nil {
		return "", err
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		return "", ErrTokenUnavailable
	}
	next := StoredToken{
		ShopID: shopID, AccessToken: refreshed.AccessToken, RefreshToken: refreshed.RefreshToken,
		ExpiresAt: m.now().Add(time.Duration(refreshed.ExpireIn) * time.Second),
	}
	if err := m.store.Save(ctx, next); err != nil {
		return "", err
	}
	return next.AccessToken, nil
}
