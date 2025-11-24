package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
)

func main() {
	ctx := context.Background()
	// 1. 定义 Edge 浏览器的路径
	edgePath := `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`
	// 注意：在 Windows 字符串中使用反引号 ` 或双反斜杠 \\ 来避免转义问题。

	// 2. 配置 Config 结构体
	config := &browseruse.Config{
		// 将 Edge 路径赋值给 ChromeInstancePath
		ChromeInstancePath: edgePath,
		// 建议设置为 Headless: true，除非您需要看到浏览器窗口
		Headless: false,
	}
	but, err := browseruse.NewBrowserUseTool(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://www.bilibili.com"
	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	time.Sleep(10 * time.Second)
	but.Cleanup()
}
