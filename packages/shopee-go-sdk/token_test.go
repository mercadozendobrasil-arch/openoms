package shopee

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type memoryTokenStore struct{ token StoredToken }
func (s *memoryTokenStore) Load(context.Context,string)(StoredToken,error){return s.token,nil}
func (s *memoryTokenStore) Save(_ context.Context,t StoredToken)error{s.token=t;return nil}

func TestTokenManagerKeepsFreshToken(t *testing.T){
	now:=time.Unix(1700000000,0)
	c:=NewClient(Config{PartnerID:1,PartnerKey:"key",Now:func()time.Time{return now}})
	store:=&memoryTokenStore{token:StoredToken{ShopID:"123",AccessToken:"fresh",RefreshToken:"refresh",ExpiresAt:now.Add(time.Hour)}}
	m:=NewTokenManager(c,store,5*time.Minute)
	got,err:=m.AccessToken(context.Background(),"123")
	if err!=nil||got!="fresh"{t.Fatalf("got=%q err=%v",got,err)}
}

func TestTokenManagerRefreshesExpiringToken(t *testing.T){
	now:=time.Unix(1700000000,0)
	srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if r.URL.Path!="/api/v2/auth/access_token/get"{t.Fatalf("path=%s",r.URL.Path)}
		_,_=w.Write([]byte(`{"access_token":"new-access","refresh_token":"new-refresh","expire_in":14400,"error":"","message":""}`))
	}))
	defer srv.Close()
	c:=NewClient(Config{BaseURL:srv.URL,PartnerID:1,PartnerKey:"key",HTTPClient:srv.Client(),Now:func()time.Time{return now}})
	store:=&memoryTokenStore{token:StoredToken{ShopID:"123",AccessToken:"old",RefreshToken:"refresh",ExpiresAt:now.Add(time.Minute)}}
	m:=NewTokenManager(c,store,5*time.Minute)
	got,err:=m.AccessToken(context.Background(),"123")
	if err!=nil||got!="new-access"{t.Fatalf("got=%q err=%v",got,err)}
	if store.token.RefreshToken!="new-refresh"||!store.token.ExpiresAt.Equal(now.Add(4*time.Hour)){t.Fatalf("stored=%+v",store.token)}
}
