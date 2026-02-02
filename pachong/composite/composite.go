package composite

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func Composite() {
	url := launcher.New().
		Headless(false).
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		MustLaunch()
	//启动
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://sale.1688.com/factory/u0vjcc4j.html?spm=a260k.home2025.centralDoor.ddoor.66333597BBbHgE&topOfferIds=1005591171200")
	page.MustWaitLoad()
	page.MustWaitIdle()
	err := page.Timeout(10 * time.Second).MustElement(".offerItem").WaitVisible()
	if err != nil {
		panic("获取商品列表超时: " + err.Error())
	}

	page.Timeout(10 * time.Second)
	err = page.Mouse.Scroll(0, 1000, 2)
	if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(30 * time.Second)
	err = page.Mouse.Scroll(0, 500, 0)
	if err != nil {
		log.Fatalln(err)
	}
	page.MustWaitLoad()

	elements := page.MustElements(".offerItem")
	fmt.Println("商品列表数量:", len(elements))
	if len(elements) == 0 {
		panic("商品列表为空")
	}

	for i, element := range elements {
		err = element.WaitVisible()
		if err != nil {
			panic("商品列表第" + fmt.Sprintf("%d", i) + "个商品超时: " + err.Error())
		}

		mustElements := element.MustElements(".text")

		var priceParts []string

		for _, mustElement := range mustElements {
			text := mustElement.MustText()
			if strings.Contains(text, "¥") || strings.Contains(text, ".") ||
				len(text) > 0 && text[0] >= '0' && text[0] <= '9' {
				priceParts = append(priceParts, text)
			}
		}

		price := strings.Join(priceParts, "")
		if price == "" {
			priceBox, err := element.Element("div[class*='text'], div[class*='text']")
			if err != nil {
				price, _ = priceBox.Text()
			}
		}
		title := element.MustElement(".offerTitle")
		titleText := title.MustText()

		url := element.MustElement(".offerImg")
		attribute, _ := url.Attribute("src")

		fmt.Printf("商品列表第%d个商品价格: %s, 商品名称: %s, 商品链接:https:%s \n", i+1, price, titleText, *attribute)
	}
}
