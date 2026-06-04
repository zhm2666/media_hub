package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// UserServiceClient 用户中心服务客户端
type UserServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// UserInfo 用户信息结构
type UserInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// NewUserServiceClient 创建用户服务客户端
func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyToken 验证 JWT Token 并获取用户信息
// 调用 user 服务的 /api/v1/login/check/auth 接口
func (c *UserServiceClient) VerifyToken(token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/api/v1/login/check/auth?access_token=%s", c.baseURL, token)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求用户中心失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("验证失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Code    int       `json:"code"`
		Message string    `json:"message"`
		Data    *UserInfo `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("验证失败: %s", result.Message)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("用户信息为空")
	}

	return result.Data, nil
}

// GetLoginURL 获取登录跳转 URL
func (c *UserServiceClient) GetLoginURL(sys string) string {
	return fmt.Sprintf("%s?sys=%s", c.baseURL, sys)
}
