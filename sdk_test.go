package golang_sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestNewDigitalSignSDK(t *testing.T) {
	test, _ := ioutil.ReadFile("test.key.json")
	var user struct{
		Key string `json:"key"`
		Secret string `json:"secret"`
	}
	_ = json.Unmarshal(test, &user)
	var verifyFn VerifyFunc = func(m map[string]interface{}) bool {
		return true
	}
	sdk := NewDigitalSignSDK(user.Key, user.Secret, "http://example.com", verifyFn)
	rsp, err := sdk.NativeOrder("testorder-123", "测试商品", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if rsp.Success{
		t.Log(rsp.Data.RequestData.CodeURL)
	} else {
		t.Log(rsp.Errors)
	}
}
func ExampleNewDigitalSignSDK() {
	test, _ := ioutil.ReadFile("test.key.json")
	var user struct{
		Key string `json:"key"`
		Secret string `json:"secret"`
	}
	_ = json.Unmarshal(test, &user)
	fmt.Println(user.Key, user.Secret)
	var verifyFn VerifyFunc = func(m map[string]interface{}) bool {
		return true
	}
	sdk := NewDigitalSignSDK(user.Key, user.Secret, "https://example.com", verifyFn)
	rsp, err := sdk.MiniAppOrder("testorder123", "测试商品", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}