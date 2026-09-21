package shopee

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestAuthorizationURL(t *testing.T) {
	c := NewClient(Config{BaseURL: "https://example.test", PartnerID: 1001, PartnerKey: "secret", APIVersion: "v2"})
	got, err := c.AuthorizationURL("https://app.test/callback", "state-1", 1700000000)
	if err != nil { t.Fatal(err) }
	u, _ := url.Parse(got)
	if u.Path != "/api/v2/shop/auth_partner" { t.Fatalf("path=%s", u.Path) }
	if u.Query().Get("partner_id") != "1001" || u.Query().Get("sign") == "" { t.Fatalf("bad query: %s", u.RawQuery) }
}

func TestGetAccessTokenAcceptsTopLevelPayload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/auth/token/get" { t.Fatalf("path=%s", r.URL.Path) }
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil { t.Fatal(err) }
		if body["partner_id"] != float64(1001) || body["code"] != "code" || body["shop_id"] != float64(123) { t.Fatalf("body=%v", body) }
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"access","refresh_token":"refresh","expire_in":14400,"request_id":"req-1","shop_id_list":[123],"error":"","message":""}`))
	}))
	defer srv.Close()

	c := NewClient(Config{
		BaseURL: srv.URL, PartnerID: 1001, PartnerKey: "secret", APIVersion: "v2",
		HTTPClient: srv.Client(), Now: func() time.Time { return time.Unix(1700000000, 0) },
	})
	token, err := c.GetAccessToken(context.Background(), "code", "123")
	if err != nil { t.Fatal(err) }
	if token.AccessToken != "access" || token.RefreshToken != "refresh" || token.ExpireIn != 14400 { t.Fatalf("token=%+v", token) }
}

func TestSignedBusinessRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/shop/get_shop_info" { t.Fatalf("path=%s", r.URL.Path) }
		q := r.URL.Query()
		if q.Get("partner_id") != "1001" || q.Get("shop_id") != "123" || q.Get("access_token") != "token" || q.Get("sign") == "" { t.Fatalf("query=%v", q) }
		_, _ = w.Write([]byte(`{"response":{"shop_id":123,"shop_name":"HAMI"},"request_id":"req"}`))
	}))
	defer srv.Close()
	c := NewClient(Config{BaseURL:srv.URL, PartnerID:1001, PartnerKey:"secret", HTTPClient:srv.Client(), Now:func() time.Time{return time.Unix(1700000000,0)}})
	var out struct { ShopID int64 `json:"shop_id"`; ShopName string `json:"shop_name"` }
	if err:=c.Do(context.Background(), Request{Method:http.MethodGet, Path:"/shop/get_shop_info", AccessToken:"token", ShopID:"123"}, &out); err!=nil {t.Fatal(err)}
	if out.ShopID!=123 || out.ShopName!="HAMI" {t.Fatalf("out=%+v",out)}
}
