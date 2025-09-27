package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// setupRouter 配置 Gin 路由
func setupRouter() *gin.Engine {
	r := gin.Default()

	// 为上传文件创建一个存储目录
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		panic(fmt.Sprintf("创建上传目录失败: %v", err))
	}

	// 配置一个静态文件路由，用于存放上传的文件
	// 这样用户可以通过 /uploads/filename 来访问上传的文件
	r.Static("/uploads", "./uploads")

	// 提供一个简单的 HTML 表单用于测试
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Gin 文件上传示例</title>
		</head>
		<body>
			<h2>单文件上传</h2>
			<form action="/upload/single" method="post" enctype="multipart/form-data">
				<input type="file" name="upload_file">
				<input type="submit" value="上传文件">
			</form>

			<hr>

			<h2>多文件上传</h2>
			<form action="/upload/multiple" method="post" enctype="multipart/form-data">
				<input type="file" name="upload_files" multiple>
				<input type="submit" value="上传多个文件">
			</form>
		</body>
		</html>
		`))
	})

	// 1. 处理单文件上传的路由
	r.POST("/upload/single", func(c *gin.Context) {
		// "upload_file" 是 HTML 表单中 <input type="file"> 的 name 属性值
		file, err := c.FormFile("upload_file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 定义文件保存的路径
		// gin.H 是 map[string]interface{} 的一个快捷方式
		dst := fmt.Sprintf("uploads/%s", file.Filename)

		// 使用 c.SaveUploadedFile 保存文件到指定目录
		// 这个方法封装了文件创建和 io.Copy 的操作，非常方便
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"message":  fmt.Sprintf("文件 '%s' 上传成功!", file.Filename),
			"filepath": dst,
		})
	})

	// 2. 处理多文件上传的路由
	r.POST("/upload/multiple", func(c *gin.Context) {
		// 使用 c.MultipartForm() 获取多文件上传的表单
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// "upload_files" 是 HTML 表单中 <input type="file"> 的 name 属性值
		files := form.File["upload_files"]

		var uploadedFiles []string
		// 遍历所有上传的文件
		for _, file := range files {
			dst := fmt.Sprintf("uploads/%s", file.Filename)
			// 逐个保存文件
			if err := c.SaveUploadedFile(file, dst); err != nil {
				// 如果其中一个文件保存失败，可以根据业务需求选择继续或中断
				// 这里我们选择返回错误并中断
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("文件 '%s' 保存失败: %v", file.Filename, err),
				})
				return
			}
			uploadedFiles = append(uploadedFiles, file.Filename)
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"message":        "所有文件上传成功!",
			"uploaded_files": uploadedFiles,
		})
	})

	return r
}

func main() {
	r := setupRouter()

	fmt.Println("服务器正在启动，请打开浏览器访问 http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("启动服务器失败: %v\n", err)
	}
}
