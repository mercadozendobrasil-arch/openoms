package shopee

import (
	"testing"
	"github.com/google/uuid"
)

func TestBuildCatalogMappingPreservesPlatformIDs(t *testing.T){
	m:=BuildCatalogMapping(uuid.New(),uuid.New(),ListingSnapshot{ShopID:"100",ItemID:"200",ItemName:"Capa",ItemSKU:"ITEM-SKU",Status:"NORMAL",Models:[]ModelSnapshot{{ModelID:"300",ModelSKU:"TS-Rosa-A9plus",ModelName:"Rosa"}}})
	if m.Product.Source!="shopee"||m.Product.ExternalID==nil||*m.Product.ExternalID!="200"{t.Fatalf("product=%+v",m.Product)}
	if m.Listing.Status!="active"||m.Listing.ExternalID==nil||*m.Listing.ExternalID!="200"{t.Fatalf("listing=%+v",m.Listing)}
	if len(m.Variants)!=1||m.Variants[0].SKU==nil||*m.Variants[0].SKU!="TS-Rosa-A9plus"{t.Fatalf("variants=%+v",m.Variants)}
}
