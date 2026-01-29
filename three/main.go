package main

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	clickE()
}

func clickE() {
	url := launcher.New().
		Headless(false). // 调试阶段设为false，确认成功后再改为true
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		Set("disable-blink-features", "AutomationControlled"). // 隐藏自动化特征
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// 2. 访问列表页（移除URL中的空格）
	listPage := browser.MustPage("https://air.1688.com/kapp/channel-fe/cps-4c-pc/sytm?type=1&offerIds=660390230106,574965204819,949033739317")
	listPage.MustWaitLoad()

	// 3. 等待商品列表加载完成
	// 先等待列表容器出现
	listPage.Timeout(10 * time.Second).MustElement(`.offer-list.offer-list-layout`)

	// 再等待具体的商品项出现
	err := listPage.Timeout(10 * time.Second).MustElement(`.offer-item`).WaitVisible()
	if err != nil {
		fmt.Println("等待商品项超时:", err)
		// 截图调试
		listPage.MustScreenshot("list_debug.png")
		return
	}

	// 4. 简短暂停，确保所有商品加载完成
	time.Sleep(2 * time.Second)

	// 5. 获取所有商品项（超链接）
	offerItems := listPage.MustElements(`.offer-item`)
	fmt.Printf("找到 %d 个商品项\n", len(offerItems))

	if len(offerItems) == 0 {
		fmt.Println("未找到任何商品项，请检查选择器")
		return
	}

	// 6. 点击第一个商品
	fmt.Println("正在点击第一个商品...")
	firstItem := offerItems[0]

	// 滚动到元素位置确保可点击
	firstItem.MustScrollIntoView()
	time.Sleep(500 * time.Millisecond)

	// 点击元素
	firstItem.MustClick()

	// 7. 等待详情页加载完成
	// 注意：点击后可能会打开新标签页，或者在本页跳转
	time.Sleep(3 * time.Second)

	// 获取当前页面信息（处理可能的页面跳转情况）
	var detailPage *rod.Page

	// 检查是否有新窗口打开
	pages := browser.MustPages()
	if len(pages) > 1 {
		// 如果开了新标签页，切换到最新标签页
		detailPage = pages[len(pages)-1]
	} else {
		// 否则使用当前页（本页跳转）
		detailPage = listPage
	}

	// 8. 确认详情页已加载
	detailPage.MustWaitLoad()

	// 等待详情页关键元素出现（例如商品标题或图片）
	detailPage.Timeout(10 * time.Second).MustElement(`[class*="title"], [class*="detail"], img`)

	fmt.Printf("已进入详情页，URL: %s\n", detailPage.MustInfo().URL)

	// 9. 在详情页执行操作（示例：获取标题和截图）
	title, err := detailPage.Element(`h1, .title, [class*="title"]`)
	if err == nil {
		titleText, _ := title.Text()
		fmt.Printf("商品标题: %s\n", titleText)
	}

	// 截图保存详情页
	detailPage.MustScreenshot("detail_page.png")
	fmt.Println("详情页已截图保存为 detail_page.png")

	// 可选：获取详情页的特定信息
	// 例如获取所有图片
	images := detailPage.MustElements(`img`)
	fmt.Printf("详情页共有 %d 张图片\n", len(images))
	for i, img := range images {
		if src, err := img.Attribute("src"); err == nil && src != nil {
			fmt.Printf("图片 %d: %s\n", i+1, *src)
		}
	}
}

func getSrd() {
	// 启动浏览器（调试时可设为false查看过程）
	url := launcher.New().
		Headless(false). // 先设为false调试，确认能抓取到后再改回true
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// 访问1688页面（移除URL中的空格）
	page := browser.MustPage("https://air.1688.com/kapp/channel-fe/cps-4c-pc/sytm?type=1&offerIds=660390230106,574965204819,949033739317")

	// 等待页面基本加载完成
	page.MustWaitLoad()

	// 关键：等待特定元素出现（修正选择器并增加等待时间）
	err := page.Timeout(10 * time.Second).MustElement(`.offer-item__img`).WaitVisible()
	if err != nil {
		fmt.Println("等待元素超时:", err)
		return
	}

	// 简短暂停，确保所有图片加载完成
	time.Sleep(2 * time.Second)

	// 获取所有匹配的元素（修正CSS选择器）
	elements := page.MustElements(`.offer-item__img`)

	fmt.Printf("找到 %d 个商品图片:\n", len(elements))

	// 遍历并提取src属性
	for i, el := range elements {
		// 确保元素可见
		if el.MustVisible() {
			src, err := el.Attribute("src")
			if err == nil && src != nil {
				fmt.Printf("%d. %s\n", i+1, *src)
			}
		}
	}
}
