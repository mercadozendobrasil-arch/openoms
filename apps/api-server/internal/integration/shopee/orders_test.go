package shopee
import"testing"
func TestOrderStatusMappingInput(t *testing.T){s:=OrderSnapshot{OrderSN:"A",Status:"READY_TO_SHIP",Items:[]OrderItemSnapshot{{ItemID:"1",ModelID:"2",ModelSKU:"TS-Rosa-A9plus",Quantity:1}}};if s.OrderSN!="A"||s.Items[0].ModelSKU!="TS-Rosa-A9plus"{t.Fatal(s)}}
