package shopee

import (
	"encoding/json"
	"fmt"
	"time"
)

func RefreshedCredentialJSON(existing []byte, accessToken, refreshToken string, expiry time.Time) ([]byte, error) {
	if accessToken == "" || refreshToken == "" {
		return nil, fmt.Errorf("shopee: refresh response missing access_token or refresh_token")
	}
	var creds Credentials
	if err := json.Unmarshal(existing, &creds); err != nil {
		return nil, fmt.Errorf("shopee: parse credentials: %w", err)
	}
	creds.AccessToken = accessToken
	creds.RefreshToken = refreshToken
	creds.TokenExpiry = expiry.UTC().Format(time.RFC3339)
	return json.Marshal(creds)
}

func ReconnectCredentialJSON(partnerID int64, partnerKey, shopID, accessToken, refreshToken string, expiry time.Time, sandbox bool) ([]byte, error) {
	if partnerID == 0 || partnerKey == "" || shopID == "" || accessToken == "" || refreshToken == "" {
		return nil, fmt.Errorf("shopee: incomplete reconnect credentials")
	}
	return json.Marshal(Credentials{
		PartnerID: partnerID, PartnerKey: partnerKey, ShopID: shopID,
		AccessToken: accessToken, RefreshToken: refreshToken,
		TokenExpiry: expiry.UTC().Format(time.RFC3339), Sandbox: sandbox,
	})
}
