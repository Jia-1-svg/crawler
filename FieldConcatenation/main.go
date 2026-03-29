package main

import (
	"bytes"
	"fmt"
	"strings"
)

// 1. 使用 + 号拼接（函数1）
func ConcatWithPlus(a, b, c string) string {
	return a + b + c
}

// 2. 使用 strings.Builder 拼接（函数2 —— 性能最优）
func ConcatWithBuilder(a, b, c string) string {
	var builder strings.Builder
	builder.WriteString(a)
	builder.WriteString(b)
	builder.WriteString(c)
	return builder.String()
}

// 3. 使用 bytes.Buffer 拼接（函数3）
func ConcatWithBuffer(a, b, c string) string {
	var buf bytes.Buffer
	buf.WriteString(a)
	buf.WriteString(b)
	buf.WriteString(c)
	return buf.String()
}

func main() {
	// 测试数据
	s1 := "Hello "
	s2 := "Go "
	s3 := "World"

	// 调用3种拼接方法
	res1 := ConcatWithPlus(s1, s2, s3)
	res2 := ConcatWithBuilder(s1, s2, s3)
	res3 := ConcatWithBuffer(s1, s2, s3)

	fmt.Println("+  号拼接结果：", res1)
	fmt.Println("Builder 拼接结果：", res2)
	fmt.Println("Buffer 拼接结果：", res3)
}
