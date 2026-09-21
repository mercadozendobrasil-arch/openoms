package shopee

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetShopInfo(t *testing.T){
	srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if r.URL.Path!="/api/v2/shop/get_shop_info"{t.Fatalf("path=%s",r.URL.Path)}
		if r.URL.Query().Get("shop_id")!="123"||r.URL.Query().Get("access_token")!="token"{t.Fatalf("query=%v",r.URL.Query())}
		_,_=w.Write([]byte(`{"response":{"shop_id":123,"shop_name":"HAMI","region":"BR"},"request_id":"r1"}`))
	}))
	defer srv.Close()
	c:=NewClient(Config{BaseURL:srv.URL,PartnerID:1,PartnerKey:"key",HTTPClient:srv.Client()})
	info,err:=c.GetShopInfo(context.Background(),"token","123")
	if err!=nil{t.Fatal(err)}
	if info.ShopID!=123||info.ShopName!="HAMI"||info.Region!="BR"{t.Fatalf("info=%+v",info)}
}
