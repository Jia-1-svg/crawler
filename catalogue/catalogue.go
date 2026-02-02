package catalogue

import "os"

var str = []string{
	"api-gateway",
	"api-gateway/basic",
	"api-gateway/basic/cmd",
	"api-gateway/basic/init",
	"api-gateway/basic/config",
	"api-gateway/proto",
	"api-gateway/handler",
	"api-gateway/handler/request",
	"api-gateway/handler/response",
	"api-gateway/handler/service",
	"api-gateway/middleWare",
	"api-gateway/router",
	"api-gateway/pkg",
	"bff-sev",
	"bff-sev/basic",
	"bff-sev/basic/cmd",
	"bff-sev/basic/init",
	"bff-sev/basic/config",
	"bff-sev/proto",
	"bff-sev/handler",
	"bff-sev/model",
}

func CreateDir() {
	for _, s := range str {
		err := os.MkdirAll(s, os.ModePerm)
		if err != nil {
			panic("目录创建失败")
		}
		println("目录创建成功")
	}
}
