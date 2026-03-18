package pkg

import (
	"fmt"
	"strconv"
	"zhongyao/aa/crawler/practice/rpc/basic/config"

	"github.com/smartwalle/alipay/v3"
)

func AliPay(OutTradeNo string, TotalAmount float64) string {
	aliConfig := config.Config.ALiPayConfig
	var privateKey = aliConfig.Key // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	appId := aliConfig.Id
	client, err := alipay.New(appId, privateKey, false)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = "https://51147577.r8.vip.cpolar.cn/notify/pay"
	p.ReturnURL = aliConfig.ReturnUrl
	p.Subject = "付款啊啊ε＝ε＝ε＝(#>д<)ﾉ"
	p.OutTradeNo = OutTradeNo
	p.TotalAmount = strconv.FormatFloat(TotalAmount, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	fmt.Println(payURL)
	return payURL
}
