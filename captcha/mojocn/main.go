package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

// 验证码存储在内存中（也可以换成 Redis）
var store = base64Captcha.DefaultMemStore

// 生成验证码
func generateCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80) // 高度80, 宽度240, 5位数字
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, ans, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证码生成失败"})
		return
	}
	fmt.Printf("验证码 id:=%s ans=%s\n", id, ans)

	c.JSON(http.StatusOK, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s, // Base64 编码的图片
	})
}

// 校验验证码
func verifyCaptcha(c *gin.Context) {
	var req struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if store.Verify(req.ID, req.Value, true) { // true 表示验证成功后清除
		fmt.Printf("验证码 id:=%s value=%s\n", req.ID, req.Value)
		c.JSON(http.StatusOK, gin.H{"message": "验证成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "验证码错误"})
	}
}

func main() {
	r := gin.Default()
	r.GET("/captcha", generateCaptcha)
	r.POST("/verify", verifyCaptcha)
	r.Run(":8080")
}
