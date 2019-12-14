# 山城清算 Golang SDK
## 安装
`go get https://github.com/023pay/golang-sdk`
## 定义

## 使用
- 实现VerifyFunc验证异步回调结果函数
    
    `type VerifyFunc func(data map[string]interface{}) bool`
    
    VerifyFunc与sdk.NotifyHandler相关。
    支付完成后，微信会把相关支付结果及用户信息通过数据流的形式发送给商户，商户需要接收处理，并按文档规范返回应答。
    您可以使用我们包装好的`sdk.NotifyHandler`来进行处理，或您可以自行编写处理函数。
    
    若您使用`sdk.NotifyHandler`则需要配置`VerifyFunc`，sdk中自带默认返回true的`DefaultVerifyFunc`，参考如下：
    ```go
    var DefaultVerifyFunc VerifyFunc = func(data map[string]interface{}) bool {
    		return true
    }
    ```
    您也可以自定义验证函数并注册到sdk中
    ```go
    var CustomVerifyFunc VerifyFunc = func(data map[string]interface{}) bool {
      // 检查订单号
      if tradeID, exist := data["out_trade_no"]; !exist{
          return false
      }   
      // 验证订单是否存在于您的数据库中
      if !Database.Order.Exist(tradeID){
          return false
      }
      // 更多验证内容
      ...
      return true
    }
    ```
    根据异步回调数据验证支付结果，微信返回的异步回调数据参考: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_7，您也可以在测试时打印`VerifyFunc`参数来了解具体参数内容。
  
- 创建SDK实例
    ```go
    key := "EzUhuxkSZAlTBqeK"
    secret := "Fadjzvg6Nzskz145DfQoEpSHCnlR9FD7"
    notifyURL := "http://example.com/callback/wechat"
    verifyFn := golang_sdk.DefaultVerifyFunc
    sdk := NewDigitalSignSDK("key", "secret", notifyURL, verifyFn)
    ```

- Native下单
    ```go
    var orderID = "0HQ48HEM3SSE"
    var amount = 1501 // 15.01元
    rsp, err := sdk.NativeOrder(orderID, "测试商品", amount)
    if err != nil {
    	// 出现网络错误
        fmt.Println(err)
        return
    }
    ```
- MiniApp下单(未实现)
    ```go
    var orderID = "0HQ48HEM3SSE"
    var amount = 1501 // 15.01元
    rsp, err := sdk.MiniAppOrder(orderID, "测试商品", amount)
    if err != nil {
    	// 出现网络错误
        fmt.Println(err)
        return
    }
    ```
- H5下单(未实现)
    ```go
    var orderID = "0HQ48HEM3SSE"
    var amount = 1501 // 15.01元
    rsp, err := sdk.H5Order(orderID, "测试商品", amount)
    if err != nil {
    	// 出现网络错误
        fmt.Println(err)
        return
    }
   ```
- 下单结果处理
    ```go
    // 创建成功
    if rsp.Success{
        // 打印二维码
        fmt.Println(rsp.Data.RequestData.CodeURL)
    } else { // 创建失败，有可能是参数不合法
        t.Log(rsp.Errors) // 打印错误
    }
   ```
- 异步通知回调处理

  为了方便处理异步回调结果，我们封装了http.HandleFunc，您需要配置该函数的外网访问路径和验证函数，随后您可以使用以下demo来进行测试
  
  ```go
  http.HandleFunc("/callback/wechat", sdk.CallbackHandler)
  http.ListenAndServe(":8080", nil)
  ```