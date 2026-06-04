package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"demo-server/utils"
)

type UserCenterHandler struct {
	userCenterURL string
	sysName       string
}

func NewUserCenterHandler(url, sysName string) *UserCenterHandler {
	return &UserCenterHandler{
		userCenterURL: url,
		sysName:       sysName,
	}
}

type LoginMethodsResponse struct {
	Gitlab   string    `json:"gitlab"`
	WxQrcode *WxQrcode `json:"wx_qrcode"`
}

type WxQrcode struct {
	Ticket    string `json:"ticket"`
	QrCodeUrl string `json:"qr_code_url"`
}

type WxCallbackResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

type CheckLoginResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

func (h *UserCenterHandler) GetLoginMethods(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/login/methods?sys=%s", h.userCenterURL, h.sysName)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取登录方式失败"})
		return
	}
	defer resp.Body.Close()

	var result LoginMethodsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析响应失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserCenterHandler) GetWxQrcode(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/login/methods?sys=%s", h.userCenterURL, h.sysName)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取二维码失败"})
		return
	}
	defer resp.Body.Close()

	var result LoginMethodsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析响应失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ticket":     result.WxQrcode.Ticket,
		"qrcode_url": result.WxQrcode.QrCodeUrl,
	})
}

func (h *UserCenterHandler) CheckWxScanStatus(c *gin.Context) {
	ticket := c.Query("ticket")
	if ticket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少ticket参数"})
		return
	}

	url := fmt.Sprintf("%s/api/v1/login/official/callback?ticket=%s&sys=%s", h.userCenterURL, ticket, h.sysName)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查扫码状态失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		if redirectURL, ok := result["redirect_url"].(string); ok && redirectURL != "" {
			c.JSON(http.StatusOK, gin.H{
				"status":       "scanned",
				"redirect_url": redirectURL,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "waiting",
	})
}

func (h *UserCenterHandler) CheckToken(c *gin.Context) {
	token := c.Query("access_token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少access_token"})
		return
	}

	url := fmt.Sprintf("%s/api/v1/login/check/auth?access_token=%s", h.userCenterURL, token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证Token失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result CheckLoginResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析响应失败"})
			return
		}

		demoToken, err := utils.GenerateTokenWithUserCenter(result.ID, result.Name, result.AvatarUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Token失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": demoToken,
			"user": gin.H{
				"id":       result.ID,
				"username": result.Name,
				"email":    "",
				"avatar":   result.AvatarUrl,
			},
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Token无效"})
}
