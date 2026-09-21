# Shopee SDK Reuse Audit

Source audited: `mercadozendobrasil-arch/hami-erp/apps/api/src/shopee-sdk`.

## Decision

Do not redesign Shopee integration from zero. The existing TypeScript SDK is the behavioral reference for the OpenOMS HAMI integration. OpenOMS remains Go-native, so protocol behavior and tests are ported/adapted rather than coupling the Go core to NestJS.

## Existing SDK coverage

| Area | Existing coverage | HAMI/OpenOMS action |
|---|---|---|
| Core client | signed requests, HTTP wrapper, error mapper, payload conversion | Port protocol behavior to Go |
| Auth | authorization URL, token exchange, token refresh | Port first |
| Shop | shop info, profile, profile update | Port first |
| Product | categories, attributes, brands, items, models, price, stock, violations | Port after auth |
| Orders | list/detail, notes, cancellation | Map into OpenOMS orders |
| Logistics | channels, shipping params, ship order, shipping documents | Reuse behavior |
| Payments | escrow, payout, wallet transactions | Feed finance/profit layer |
| Media | media API module exists | Port when product publishing starts |
| Invoice | invoice API module exists | Review for Brazil applicability |
| Webhook | signature verification, parsing | Adapt to OpenOMS webhook_events |

## Verified regression assets to preserve

The existing `apps/api/test/shopee-sdk.spec.ts` verifies:
- sandbox fallback and production credential safety
- Shopee V2 path construction
- snake_case payload conversion
- `/api/v2/product/add_item`
- `/api/v2/auth/token/get`
- top-level token response parsing
- rejection-sensitive order optional fields

These behaviors should become Go tests before replacing or retiring the TypeScript implementation.

## OpenOMS mapping

```
Shopee shop              -> integrations + HAMI shop registry
Shopee item              -> product_listings
Shopee model/SKU         -> product variant / Master SKU mapping
Shopee order             -> orders
Shopee stock             -> warehouse_stock synchronization boundary
Shopee webhook           -> webhook_events -> worker reconciliation
Shopee escrow/payment    -> HAMI finance/profit inputs
Shopee Ads               -> separate Ads ingestion/profit domain
```

## Migration order

1. Signature + HTTP + error envelope
2. Auth + token refresh
3. Shop registry
4. Product/item/model read sync
5. Order read sync
6. Inventory write sync
7. Webhook ingestion + reconciliation
8. Logistics
9. Payments/escrow
10. Product publishing/write operations
11. Shopee Ads adapter

## Guardrails

- Never expose partner_key/access_token/refresh_token to the frontend.
- Keep sandbox and production configuration separate.
- Keep all synchronization idempotent.
- Preserve Shopee external IDs.
- Treat webhook as a signal; retain pull reconciliation.
- Do not guess undocumented endpoint behavior.
- Before enabling production writes, re-verify the relevant Shopee Open Platform endpoint and permission.
- Keep Shopee Ads separate from order/inventory transaction processing.
