package shopee
import ("strings";"github.com/google/uuid";"github.com/openoms-org/openoms/apps/api-server/internal/model")
func BuildSKUMappings(tenantID,integrationID uuid.UUID,snapshot ListingSnapshot,mapping CatalogMapping) []model.MarketplaceSKUMapping {
 out:=make([]model.MarketplaceSKUMapping,0,len(snapshot.Models)+1)
 if len(snapshot.Models)==0 {var sku *string;if v:=strings.TrimSpace(snapshot.ItemSKU);v!=""{sku=&v};return append(out,model.MarketplaceSKUMapping{ID:uuid.New(),TenantID:tenantID,IntegrationID:integrationID,ProductID:mapping.Product.ID,Provider:"shopee",ExternalItemID:snapshot.ItemID,ExternalSKU:sku,MappingSource:"import",Active:true})}
 for i,m:=range snapshot.Models {if i>=len(mapping.Variants){break};var sku *string;if v:=strings.TrimSpace(m.ModelSKU);v!=""{sku=&v};variantID:=mapping.Variants[i].ID;out=append(out,model.MarketplaceSKUMapping{ID:uuid.New(),TenantID:tenantID,IntegrationID:integrationID,ProductID:mapping.Product.ID,VariantID:&variantID,Provider:"shopee",ExternalItemID:snapshot.ItemID,ExternalModelID:m.ModelID,ExternalSKU:sku,MappingSource:"import",Active:true})}
 return out
}
