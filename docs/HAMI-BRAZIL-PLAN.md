# HAMI Brazil Adaptation Plan

## Decision
Use OpenOMS as the OMS/WMS core. Keep upstream-oriented core behavior intact where possible and add Brazil/channel capabilities through adapters and additive modules.

## Reuse without rewriting
- tenants / users / RBAC
- orders and order groups
- products and product variants
- product bundles
- product listings
- warehouse stock, warehouses, stocktakes and warehouse documents
- integrations and sync jobs
- automation rules
- audit log and webhook events

## HAMI semantic mapping
- OpenOMS product / variant -> HAMI Master SKU
- product_listing -> marketplace listing
- integration -> channel account / store authorization
- warehouse_stock -> physical sellable inventory
- product_bundle -> BOM / kit relationship

## Brazil phase 1
1. Shopee BR integration
2. Multi-shop authorization and token refresh
3. Shop registry and external IDs
4. Product/listing/SKU synchronization
5. Order synchronization
6. Inventory synchronization
7. Master SKU mapping across shops
8. BOM/kit stock calculation
9. Shopee Ads data ingestion
10. Profit engine and break-even ACOS/ROAS

## Integration boundary
Shopee-specific signing, OAuth/token lifecycle, API payloads and endpoint clients must remain outside the core order/inventory domain.

Proposed package:
packages/shopee-go-sdk/

Proposed application integration:
apps/api-server/internal/integrations/shopee/

The SDK is responsible for protocol/API details. The app integration translates Shopee resources into OpenOMS domain entities.

## Profit engine
Inputs:
- paid order sales
- paid units
- Master SKU cost
- commission/service fees
- advertising spend

Outputs:
- actual ASP
- COGS
- pre-ad contribution margin
- actual profit
- ACOS / ROAS
- break-even ACOS / ROAS

Unknown SKU costs must never be guessed. Manual cost overrides have highest priority. Bundle cost may be expanded only from confirmed BOM definitions.

## Development rules
- main stays aligned with upstream.
- HAMI work is committed to hami-dev.
- Prefer additive migrations and adapters over invasive core rewrites.
- Preserve external IDs and sync auditability.
- Sync jobs must be idempotent.
- Do not couple Ads API logic to OMS order processing.
