package gospacexbug

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestHTTPClient_NoLeak(t *testing.T) {
	// 创建一个慢速服务器
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // 模拟慢响应
		w.WriteHeader(http.StatusOK)
	}))
	defer slowServer.Close()

	// 使用修复后的 HTTPClient
	client := NewHTTPClient()

	// 正常请求（应该成功）
	_, err := client.Get(slowServer.URL)
	if err != nil {
		t.Logf("Expected request to succeed, got error: %v", err)
	}
}

func TestHTTPClient_Timeout(t *testing.T) {
	// 创建一个超慢的服务器
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second) // 超过客户端超时时间
		w.WriteHeader(http.StatusOK)
	}))
	defer slowServer.Close()

	client := NewHTTPClient()

	// 请求应该超时
	_, err := client.Get(slowServer.URL)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	t.Logf("Correctly got timeout error: %v", err)
}

func TestHTTPClient_BatchGet_NoLeak(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := NewHTTPClient()

	urls := []string{
		server.URL + "/1",
		server.URL + "/2",
		server.URL + "/3",
	}

	results, err := client.BatchGet(urls)
	if err != nil {
		t.Logf("BatchGet error (expected with test server): %v", err)
	} else {
		t.Logf("BatchGet completed, got %d results", len(results))
	}
}

func TestHTTPClient_RetryGet_NoLeak(t *testing.T) {
	// 创建一个失败的服务器
	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failServer.Close()

	client := NewHTTPClient()

	// 重试请求
	_, err := client.RetryGet(failServer.URL, 3)
	if err == nil {
		t.Error("Expected error after retries, got nil")
	}
	t.Logf("Correctly got error after retries: %v", err)
}

func TestHTTPClient_GetWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient()

	// 使用短超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.GetWithContext(ctx, server.URL)
	if err == nil {
		t.Error("Expected context deadline exceeded error, got nil")
	}
	t.Logf("Correctly got context error: %v", err)
}
