package shopee

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	shopeesdk "github.com/openoms-org/openoms/packages/shopee-go-sdk"
)

// CredentialTokenStore adapts OpenOMS encrypted integration credentials to the SDK TokenStore.
// The caller supplies read/write functions so encryption and tenant isolation stay in OpenOMS.
type CredentialTokenStore struct {
	Read  func(context.Context) ([]byte, error)
	Write func(context.Context, []byte) error
}

func (s CredentialTokenStore) Load(ctx context.Context, shopID string) (shopeesdk.StoredToken, error) {
	raw, err := s.Read(ctx)
	if err != nil { return shopeesdk.StoredToken{}, err }
	var creds Credentials
	if err := json.Unmarshal(raw, &creds); err != nil {
		return shopeesdk.StoredToken{}, fmt.Errorf("shopee: parse stored credentials: %w", err)
	}
	if creds.ShopID != shopID {
		return shopeesdk.StoredToken{}, fmt.Errorf("shopee: stored shop_id mismatch")
	}
	expiry, err := time.Parse(time.RFC3339, creds.TokenExpiry)
	if err != nil { return shopeesdk.StoredToken{}, fmt.Errorf("shopee: parse token expiry: %w", err) }
	return shopeesdk.StoredToken{ShopID: creds.ShopID, AccessToken: creds.AccessToken, RefreshToken: creds.RefreshToken, ExpiresAt: expiry}, nil
}

func (s CredentialTokenStore) Save(ctx context.Context, token shopeesdk.StoredToken) error {
	raw, err := s.Read(ctx)
	if err != nil { return err }
	next, err := RefreshedCredentialJSON(raw, token.AccessToken, token.RefreshToken, token.ExpiresAt)
	if err != nil { return err }
	return s.Write(ctx, next)
}
