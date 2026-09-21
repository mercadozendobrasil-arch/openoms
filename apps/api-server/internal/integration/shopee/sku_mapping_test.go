package shopee
import ("testing";"github.com/google/uuid")
func TestBuildSKUMappingsUsesModelIdentity(t *testing.T){tenant,integration:=uuid.New(),uuid.New();s:=ListingSnapshot{ShopID:"1",ItemID:"2",Models:[]ModelSnapshot{{ModelID:"3",ModelSKU:"TS-Rosa-A9plus"}}};c:=BuildCatalogMapping(tenant,integration,s);m:=BuildSKUMappings(tenant,integration,s,c);if len(m)!=1||m[0].ExternalModelID!="3"||m[0].ExternalSKU==nil||*m[0].ExternalSKU!="TS-Rosa-A9plus"||m[0].VariantID==nil{t.Fatalf("mapping=%+v",m)}}
