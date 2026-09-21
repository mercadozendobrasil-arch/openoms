package shopee

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProductReadEndpoints(t *testing.T){
	srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		switch r.URL.Path {
		case "/api/v2/product/get_item_list":
			_,_=w.Write([]byte(`{"response":{"item":[{"item_id":11,"item_status":"NORMAL"}],"has_next_page":false}}`))
		case "/api/v2/product/get_item_base_info":
			if r.URL.Query().Get("item_id_list")!="11"{t.Fatalf("ids=%s",r.URL.Query().Get("item_id_list"))}
			_,_=w.Write([]byte(`{"response":{"item_list":[{"item_id":11,"item_name":"Case","item_sku":"TS-Rosa-A9plus"}]}}`))
		case "/api/v2/product/get_model_list":
			_,_=w.Write([]byte(`{"response":{"model":[{"model_id":21,"model_sku":"TS-Rosa-A9plus","model_name":"Rosa"}]}}`))
		default:t.Fatalf("path=%s",r.URL.Path)
		}
	})); defer srv.Close()
	c:=NewClient(Config{BaseURL:srv.URL,PartnerID:1,PartnerKey:"key",HTTPClient:srv.Client()})
	list,err:=c.GetItemList(context.Background(),"token","123","NORMAL",0,50);if err!=nil||len(list.Item)!=1{t.Fatalf("list=%+v err=%v",list,err)}
	info,err:=c.GetItemBaseInfo(context.Background(),"token","123",[]int64{11});if err!=nil||info[0].ItemSKU!="TS-Rosa-A9plus"{t.Fatalf("info=%+v err=%v",info,err)}
	models,err:=c.GetModelList(context.Background(),"token","123",11);if err!=nil||models[0].ModelID!=21{t.Fatalf("models=%+v err=%v",models,err)}
}
