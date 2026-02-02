package main

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// 图片

func main() {
	// 启动浏览器（调试时可设为false查看过程）
	url := launcher.New().
		Headless(false). // 先设为false调试，确认能抓取到后再改回true
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// 访问1688页面（移除URL中的空格）
	page := browser.MustPage("https://sale.1688.com/factory/u0vjcc4j.html?spm=a260k.home2025.centralDoor.ddoor.66333597BBbHgE&topOfferIds=1005591171200")

	// 等待页面基本加载完成
	page.MustWaitLoad()

	// 关键：等待特定元素出现（修正选择器并增加等待时间）
	err := page.Timeout(10 * time.Second).MustElement(`.offerImg`).WaitVisible()
	if err != nil {
		fmt.Println("等待元素超时:", err)
		return
	}

	// 简短暂停，确保所有图片加载完成
	time.Sleep(5 * time.Second)
	scrollTimes := 2                  // 滚动次数，按需调整（1688一般10-20次足够加载全部）
	scrollInterval := 2 * time.Second // 每次滚动后的休眠时间（给懒加载留时间，网络慢可改3s）
	fmt.Printf("开始循环滚动%d次，每次滚动后休眠%v...\n", scrollTimes, scrollInterval)

	for i := 0; i < scrollTimes; i++ {
		// 1. 执行页面滚动：整屏向下滚动（模拟用户翻页，触发懒加载）
		_, err := page.Eval(`() => {
				// 每次滚动一个可视窗口的高度，适配懒加载
				window.scrollTo(0, window.scrollY + window.innerHeight);
			}`)
		if err != nil {
			fmt.Printf("第%d次滚动失败: %v\n", i+1, err)
			continue
		}
		fmt.Printf("完成第%d次滚动，等待页面加载...\n", i+1)

		// 2. 滚动后休眠：核心步骤，让JS加载新的懒加载图片/元素
		time.Sleep(scrollInterval)
	}
	// 滚动完成后，再休眠2s，确保最后一批图片完全加载
	time.Sleep(2 * time.Second)
	fmt.Println("滚动完成，开始提取所有图片地址...")
	// 获取所有匹配的元素（修正CSS选择器）
	elements := page.MustElements(`.offerImg`)
	limit := 5 // 要舍弃的最后N个数量
	if len(elements) > limit {
		elements = elements[:len(elements)-limit] // 切片截取，只保留前 总数量-5 个
	}
	fmt.Printf("找到 %d 个商品图片:\n", len(elements))

	// 遍历并提取src属性
	for i, el := range elements {
		// 确保元素可见
		if el.MustVisible() {
			src, err := el.Attribute("src")
			if err == nil && src != nil {
				fmt.Printf("%d. https:%s\n", i+1, *src)
			}
		}
	}
}
