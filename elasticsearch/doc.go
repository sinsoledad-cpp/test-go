package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"time"
)

func DocCreate() {
	user := UserModel{
		ID:       12,
		UserName: "lisi",
		//Age:       23,
		NickName:  "夜空中最亮的lisi",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		//Title:     "今天天气很不错",
	}
	indexResponse, err := ESClient.Index().Index(user.Index()).BodyJson(user).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", indexResponse)
}
func DocDelete() {

	deleteResponse, err := ESClient.Delete().Index(UserModel{}.Index()).Id("tmcqfYkBWS69Op6Q4Z0t").Refresh("true").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(deleteResponse)
}

func DocCreateBatch() {

	list := []UserModel{
		{
			//ID:        13,
			//UserName:  "lisi",
			//NickName:  "夜空中最亮的李四",
			Title:     "这是我的生活",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			//ID:        14,
			//UserName:  "zhangsan",
			//NickName:  "张三",
			Title:     "你好啊，枫枫",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			//ID:        14,
			//UserName:  "zhangsan",
			//NickName:  "张三",
			Title:     "这是我的枫枫",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	bulk := ESClient.Bulk().Index(UserModel{}.Index()).Refresh("true")
	for _, model := range list {
		req := elastic.NewBulkCreateRequest().Doc(model)
		bulk.Add(req)
	}
	res, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Succeeded())
}

func DocDeleteBatch() {
	idList := []string{
		"tGcofYkBWS69Op6QHJ2g",
		"tWcpfYkBWS69Op6Q050w",
	}
	bulk := ESClient.Bulk().Index(UserModel{}.Index()).Refresh("true")
	for _, s := range idList {
		req := elastic.NewBulkDeleteRequest().Id(s)
		bulk.Add(req)
	}
	res, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Succeeded())
}

func DocFind() {

	limit := 10
	page := 1
	from := (page - 1) * limit

	query := elastic.NewTermQuery("title.keyword", "这是我的枫枫") //精确查询 针对 keyword 字段

	//query := elastic.NewMatchQuery("title", "夜空中最亮的枫枫")
	//模糊查询 主要是査 text，也能查 keyword
	//模糊匹配 keyword 字段，是需要查完整的
	//匹配 text 字段则不用，搜完整的也会搜出很多

	//query := elastic.NewBoolQuery() //全局查询

	res, err := ESClient.Search(UserModel{}.Index()).Query(query).From(from).Size(limit).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	count := res.Hits.TotalHits.Value //总数
	fmt.Println(count)
	for _, hit := range res.Hits.Hits {
		fmt.Println(string(hit.Source))
	}
}
func DocUpdate() {
	res, err := ESClient.Update().Index(UserModel{}.Index()).Id("vmdnfYkBWS69Op6QEp2Y").
		Doc(map[string]any{
			"user_name": "你好枫枫",
		}).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", res)
}
