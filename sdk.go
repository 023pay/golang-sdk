// package golang_sdk
package golang_sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/GhostLee/golang-sdk/pkg/randstr"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// WECHAT_FAILED wechat notify handler response, you can send failed response if data is not be signed
	WECHAT_FAILED = `<xml>
<return_code><![CDATA[FAIL]]></return_code>
</xml>`
	// WECHAT_SUCCESS wechat notify handler response, you need send success response if data is valid.
	WECHAT_SUCCESS  = `<xml>
<return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
</xml>`
	host = "https://pay.digital-sign.cn"
	unifieldOrderAPI = "/api/pay/wechat"

	// payment
	native = "NATIVE"
	mweb = "MWEB"
	miniapp = "MINIAPP"
)
var	(
	ParamsInvalidErr = errors.New("params is invalid")
	AmountInvalidErr = errors.New("total amount must be greater than zero")
	DefaultVerifyFunc VerifyFunc = func(m map[string]interface{}) bool {
		return true
	}
)

// VerifyFunc verify the callback data
type VerifyFunc func(map[string]interface{}) bool
/*
sdk for digital sign wechat pay
*/

// DigitalSign sdk, must be created by NewDigitalSignSDK(key, secret, ) function
type DigitalSignSDK struct {
	secret string
	key string
	notifyUrl string
	verifyFunc VerifyFunc
	client *http.Client
	loc *time.Location
}

// NewDigitalSignSDK create a new DigitalSignSDK client, please confirm that timezone.zip in your machine, otherwise it will fatal.
func NewDigitalSignSDK(key, secret, notifyUrl string, fn VerifyFunc) *DigitalSignSDK {
	l, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalf("err load location: %v", err)
	}
	return &DigitalSignSDK{
		secret: secret,
		key: key,
		notifyUrl: notifyUrl,
		verifyFunc: fn,
		client: http.DefaultClient,
		loc: l,
	}
}
// sign sign params for authorization
func (sdk *DigitalSignSDK) sign(endpoint string, data *url.Values) {
	data.Set("nonce", randstr.String(32))
	data.Set("timestamp", time.Now().In(sdk.loc).Format("2006-01-02T15:04:05Z"))
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(sdk.secret))
	// Write Data to it
	//todo: sign get query params
	h.Write([]byte(endpoint+"?"+data.Encode()))
	// Get result and encode as base64 string
	data.Set("sign", base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func (sdk *DigitalSignSDK) NativeOrder(tradeID, title string, fee uint64) (*UnifiedOrderResponse, error) {
	data := url.Values{}
	data.Set("out_trade_no", tradeID)
	data.Set("total_fee", strconv.FormatUint( fee,10))
	data.Set("body", title)
	data.Set("accessKeyId", sdk.key)
	data.Set("trade_type", native)
	data.Set("notify_url", sdk.notifyUrl)
	return sdk.unifiedOrder(data)
}

func (sdk *DigitalSignSDK) MiniAppOrder(tradeID, title string, fee uint64) (*UnifiedOrderResponse, error) {
	data := url.Values{}
	data.Set("out_trade_no", tradeID)
	data.Set("total_fee", strconv.FormatUint( fee,10))
	data.Set("body", title)
	//data.Set("redirect_url", redirectUrl)
	data.Set("accessKeyId", sdk.key)
	data.Set("trade_type", miniapp)
	data.Set("notify_url", sdk.notifyUrl)
	return sdk.unifiedOrder(data)
}
// H5Order create h5 order
// tradeID transaction id, string, length not greater than 32 characters
func (sdk *DigitalSignSDK) H5Order(tradeID, title, redirectUrl string, fee uint64) (*UnifiedOrderResponse, error) {
	data := url.Values{}
	data.Set("out_trade_no", tradeID)
	data.Set("total_fee", strconv.FormatUint( fee,10))
	data.Set("body", title)
	data.Set("redirect_url", redirectUrl)
	data.Set("accessKeyId", sdk.key)
	data.Set("trade_type", mweb)
	data.Set("notify_url", sdk.notifyUrl)
	return sdk.unifiedOrder(data)
}

// unifiedOrder unified order implement
func (sdk *DigitalSignSDK) unifiedOrder(data url.Values) (*UnifiedOrderResponse, error) {
	sdk.sign(unifieldOrderAPI, &data)
	req, _ := http.NewRequest(http.MethodPost, host+unifieldOrderAPI, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rsp, err := sdk.client.Do(req)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	var srvRsp UnifiedOrderResponse
	err =json.Unmarshal(raw, &srvRsp)
	if err != nil {
		return nil, err
	}
	// todo: check error in response
	//if len(srvRsp.Errors)>0 {
	//	for key, val := range srvRsp.Errors{
	//
	//	}
	//}
	return &srvRsp, nil
}
// Notify callback handler
func (sdk *DigitalSignSDK) NotifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, WECHAT_FAILED)
		return
	}
	callbackRsp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, WECHAT_FAILED)
		return
	}
	defer r.Body.Close()
	data, err := XmlDecode(string(callbackRsp))
	if err != nil {
		fmt.Fprint(w, WECHAT_FAILED)
		return
	}
	if sdk.verifyFunc(data){
		fmt.Fprint(w, WECHAT_SUCCESS)
	} else {
		fmt.Fprint(w, WECHAT_FAILED)
	}

}

func XmlDecode(data string) (map[string]interface{}, error){
	decoder := xml.NewDecoder(strings.NewReader(data))
	result  := make(map[string]interface{})
	key := ""
	for{
		token, err := decoder.Token() //读取一个标签或者文本内容
		if err==io.EOF{
			return result, nil
		}
		if err!=nil{
			log.Println("parse Fail:",err)
			return result, err
		}
		switch tp := token.(type) {  //读取的TOKEN可以是以下三种类型：StartElement起始标签，EndElement结束标签，CharData文本内容
		case xml.StartElement:
			se := xml.StartElement(tp) //强制类型转换
			if se.Name.Local!="xml"{
				key=se.Name.Local
			}
			//if len(se.Attr)!=0{ //读取标签属性
			//	//fmt.Println("Attrs:",se.Attr)
			//}
			//fmt.Println("SE.NAME.SPACE:",se.Name.Space) //读取命名空间
			//fmt.Println("SE.NAME.LOCAL:",se.Name.Local) //读取标签名称
		case xml.EndElement:
			ee := xml.EndElement(tp)
			if ee.Name.Local == "xml"{
				return result, nil
			}
			//fmt.Println("EE.NAME.SPACE:",ee.Name.Space)
			//fmt.Println("EE.NAME.LOCAL:",ee.Name.Local)
		case xml.CharData: //文本数据，注意一个结束标签和另一个起始标签之间可能有空格
			cd := xml.CharData(tp)
			data := strings.TrimSpace(string(cd))
			if len(data)!=0{
				result[key] = data
			}
		}
	}
}