package shopee

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewProvider(t *testing.T) {
	raw,_:=json.Marshal(Credentials{PartnerID:1,PartnerKey:"key",ShopID:"123",AccessToken:"access",RefreshToken:"refresh",TokenExpiry:time.Now().Add(time.Hour).UTC().Format(time.RFC3339)})
	p,err:=NewProvider(raw,nil)
	if err!=nil{t.Fatal(err)}
	if p.ProviderName()!="shopee"{t.Fatalf("provider=%s",p.ProviderName())}
}

func TestRefreshedCredentialJSONPreservesShop(t *testing.T) {
	exp:=time.Date(2026,9,21,15,0,0,0,time.UTC)
	raw,_:=json.Marshal(Credentials{PartnerID:1,PartnerKey:"key",ShopID:"123",AccessToken:"old",RefreshToken:"old-r"})
	next,err:=RefreshedCredentialJSON(raw,"new","new-r",exp)
	if err!=nil{t.Fatal(err)}
	var got Credentials
	if err:=json.Unmarshal(next,&got);err!=nil{t.Fatal(err)}
	if got.ShopID!="123"||got.AccessToken!="new"||got.RefreshToken!="new-r"||got.TokenExpiry!=exp.Format(time.RFC3339){t.Fatalf("got=%+v",got)}
}
