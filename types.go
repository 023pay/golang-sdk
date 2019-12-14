package golang_sdk

type UnifiedOrderResponse struct {
	Errors map[string][]string `json:"errors"`
	Data struct {
		AccessKeyID           string `json:"access_key_id"`
		ApplymentBusinessCode int    `json:"applyment_business_code"`
		Body                  string `json:"body"`
		CreatedAt             string `json:"created_at"`
		ID                    int    `json:"id"`
		NotifyURL             string `json:"notify_url"`
		OutTradeNo            string `json:"out_trade_no"`
		RequestData           struct {
			Appid      string `json:"appid"`
			CodeURL    string `json:"code_url"`
			MchID      string `json:"mch_id"`
			NonceStr   string `json:"nonce_str"`
			PrepayID   string `json:"prepay_id"`
			ResultCode string `json:"result_code"`
			ReturnCode string `json:"return_code"`
			ReturnMsg  string `json:"return_msg"`
			Sign       string `json:"sign"`
			SubMchID   string `json:"sub_mch_id"`
			TradeType  string `json:"trade_type"`
		} `json:"request_data"`
		SubMchID  string `json:"sub_mch_id"`
		TotalFee  string `json:"total_fee"`
		TradeType string `json:"trade_type"`
		UpdatedAt string `json:"updated_at"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
} 