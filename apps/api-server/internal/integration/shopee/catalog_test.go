package shopee

import "testing"

func TestListingSnapshotIdentity(t *testing.T){
	s:=ListingSnapshot{ShopID:"123",ItemID:"11",ItemSKU:"TS-Azul-A9plus",Models:[]ModelSnapshot{{ModelID:"21",ModelSKU:"TS-Rosa-A9plus"}}}
	if s.ShopID!="123"||s.ItemID!="11"||s.Models[0].ModelSKU!="TS-Rosa-A9plus"{t.Fatalf("snapshot=%+v",s)}
}
