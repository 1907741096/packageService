package lib

import (
	"testing"
	"fmt"
	"encoding/json"
)

type BaiduClient struct {
	BaseHttp
}

func (client *BaiduClient) SetDomain(domain string) {
	client.Domain = domain
}

type Response struct {
	BaseResponse
	Data 	MerchantList
}

func (list *MerchantList) First() *Merchant { return list.List[0] }

type MerchantList struct {
	List 	[]*Merchant
}

type Merchant struct {
	AmountText 	string 		`json:"amount_text"`
	DetailUrl 	string 		`json:"detail_url"`
	LogoImgUrl 	string 		`json:"logo_img_url"`
	Name 		string 		`json:"name"`
	Id 			string 		`json:"id"`
	AuditText 	string 		`json:"audit_text"`
}

func TestCurl(t *testing.T) {
	var client BaiduClient
	var header map[string]string
	var data = map[string]string {
		"message": "hello world",
	}
	client.SetDomain("http://api-service.com")
	uri := "/v1/order/merchant-list"


	result, _ := client.Get(uri, data, header)
	var response Response
	if err := json.Unmarshal(result, &response); err != nil {
		fmt.Println("err:" ,err)
	}

	t.Log(response.Data.First().Name)
}