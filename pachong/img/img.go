package img

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// 图片
func Img(pageUrl, element string) {
	// 启动浏览器（调试时可设为false查看过程）
	url := launcher.New().
		Headless(false). // 先设为false调试，确认能抓取到后再改回true
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// 访问1688页面（移除URL中的空格）
	page := browser.MustPage(pageUrl)

	// 等待页面基本加载完成
	page.MustWaitLoad()

	// 关键：等待特定元素出现（修正选择器并增加等待时间）
	err := page.Timeout(10 * time.Second).MustElement(element).WaitVisible()
	if err != nil {
		fmt.Println("等待元素超时:", err)
		return
	}

	// 简短暂停，确保所有图片加载完成
	time.Sleep(2 * time.Second)

	// 获取所有匹配的元素（修正CSS选择器）
	elements := page.MustElements(element)

	fmt.Printf("找到 %d 个商品图片:\n", len(elements))

	// 遍历并提取src属性
	for i, el := range elements {
		// 确保元素可见
		if el.MustVisible() {
			src, err := el.Attribute("src")
			if err == nil && src != nil {
				fmt.Printf("%d. https://%s\n", i+1, *src)
			}
		}
	}
}
