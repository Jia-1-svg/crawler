package main

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	// 启动浏览器（隐藏模式，可配置用户代理）
	url := launcher.New().
		Headless(true). // 设为false可查看浏览器操作过程
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// 访问1688页面
	page := browser.MustPage("https://air.1688.com/kapp/channel-fe/cps-4c-pc/sytm?type=1&offerIds=660390230106,574965204819,949033739317")

	// 等待页面基本加载完成
	page.MustWaitLoad()

	// 1688页面有动态加载，增加额外等待时间
	// 等待元素出现，最长等待10秒
	err := page.Timeout(10 * time.Second).MustElement(`.offer-title.ellipsis`).WaitVisible()
	if err != nil {
		fmt.Println("等待元素超时:", err)
		return
	}

	// 获取所有匹配的元素
	elements := page.MustElements(`.offer-title.ellipsis`)

	fmt.Printf("找到 %d 个商品标题:\n", len(elements))

	// 遍历并提取文本
	for i, el := range elements {
		// 确保元素在DOM中且可见
		if err := el.WaitVisible(); err == nil {
			text := el.MustText()
			fmt.Printf("%d. %s\n", i+1, text)
		}
	}

	// 如果需要更精确的选择（例如只获取可见的）
	fmt.Println("\n--- 只获取可见元素 ---")
	visibleElements := page.MustElements(`.offer-title.ellipsis`)
	for i, el := range visibleElements {
		// 检查元素是否可见且在视口中
		if el.MustVisible() {
			text := el.MustText()
			fmt.Printf("%d. %s\n", i+1, text)
		}
	}
}
