package shopee

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/openoms-org/openoms/apps/api-server/internal/model"
)

// CatalogMapping is a persistence-neutral plan produced from one Shopee listing.
// It deliberately does not auto-merge unrelated products by fuzzy name.
type CatalogMapping struct {
	Product model.Product
	Variants []model.ProductVariant
	Listing model.ProductListing
}

func BuildCatalogMapping(tenantID, integrationID uuid.UUID, snapshot ListingSnapshot) CatalogMapping {
	productID:=uuid.New()
	externalID:=snapshot.ItemID
	var sku *string
	if v:=strings.TrimSpace(snapshot.ItemSKU);v!=""{sku=&v}
	meta,_:=json.Marshal(map[string]any{"shopee_shop_id":snapshot.ShopID,"shopee_item_id":snapshot.ItemID})
	p:=model.Product{ID:productID,TenantID:tenantID,ExternalID:&externalID,Source:"shopee",Name:snapshot.ItemName,SKU:sku,Metadata:meta,Images:json.RawMessage("[]"),HasVariants:len(snapshot.Models)>0}
	nowMeta,_:=json.Marshal(map[string]any{"shop_id":snapshot.ShopID,"item_id":snapshot.ItemID})
	l:=model.ProductListing{ID:uuid.New(),TenantID:tenantID,ProductID:productID,IntegrationID:integrationID,ExternalID:&externalID,Status:listingStatus(snapshot.Status),SyncStatus:"synced",StockSyncMode:"auto",Metadata:nowMeta}
	out:=CatalogMapping{Product:p,Listing:l}
	for i,m:=range snapshot.Models{
		var modelSKU *string
		if v:=strings.TrimSpace(m.ModelSKU);v!=""{modelSKU=&v}
		attrs,_:=json.Marshal(map[string]any{"shopee_model_id":m.ModelID,"shopee_model_name":m.ModelName})
		name:=m.ModelName;if strings.TrimSpace(name)==""{name=m.ModelSKU};if strings.TrimSpace(name)==""{name=m.ModelID}
		out.Variants=append(out.Variants,model.ProductVariant{ID:uuid.New(),TenantID:tenantID,ProductID:productID,SKU:modelSKU,Name:name,Attributes:attrs,Position:i,Active:true})
	}
	return out
}

func listingStatus(status string) string {
	switch strings.ToUpper(status) {
	case "NORMAL": return "active"
	case "UNLIST", "BANNED", "DELETED": return "inactive"
	default: return "inactive"
	}
}
