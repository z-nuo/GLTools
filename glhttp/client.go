package glhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

// Client 是简单的 JSON HTTP 客户端。
type Client struct {
	client *http.Client
}

// NewClient 创建带超时时间的 HTTP 客户端。
func NewClient(timeout time.Duration) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// NewClientWithHTTPClient 使用已有的 http.Client 创建客户端。
func NewClientWithHTTPClient(client *http.Client) *Client {
	return &Client{
		client: client,
	}
}

// GetJSON 发起 GET 请求并将 JSON 响应解码到 out。
func (c *Client) GetJSON(ctx context.Context, url string, headers map[string]string, out any) error {
	return c.doJSON(ctx, http.MethodGet, url, headers, nil, out)
}

// PostJSON 发起 POST 请求，将 body 编码为 JSON，并将 JSON 响应解码到 out。
func (c *Client) PostJSON(ctx context.Context, url string, headers map[string]string, body any, out any) error {
	return c.doJSON(ctx, http.MethodPost, url, headers, body, out)
}

// PostForm 发起 POST 表单请求，将 form 编码为 application/x-www-form-urlencoded，并将 JSON 响应解码到 out。
func (c *Client) PostForm(ctx context.Context, url string, headers map[string]string, form url.Values, out any) error {
	return c.doRequest(ctx, http.MethodPost, url, headers, strings.NewReader(form.Encode()), "application/x-www-form-urlencoded", out)
}

func (c *Client) doJSON(ctx context.Context, method string, url string, headers map[string]string, body any, out any) error {
	var reader io.Reader
	contentType := ""
	if body != nil {
		data, err := jsonAPI.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
		contentType = "application/json"
	}
	return c.doRequest(ctx, method, url, headers, reader, contentType, out)
}

func (c *Client) doRequest(ctx context.Context, method string, url string, headers map[string]string, reader io.Reader, contentType string, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("http status %d", resp.StatusCode)
	}
	if out == nil {
		return nil
	}
	return jsonAPI.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) httpClient() *http.Client {
	if c == nil || c.client == nil {
		return http.DefaultClient
	}
	return c.client
}
