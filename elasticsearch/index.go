package main

import (
	"context"
	"fmt"
)

// ExistsIndex 判断索引是否存在
func ExistsIndex(index string) bool {
	exists, _ := ESClient.IndexExists(index).Do(context.Background())
	return exists
}

func DeleteIndex(index string) {
	_, err := ESClient.DeleteIndex(index).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(index, "索引删除成功")
}

func CreateIndex() {
	index := "user_index"
	if ExistsIndex(index) {
		// 索引存在，先删除，在创建
		DeleteIndex(index)
	}

	createIndex, err := ESClient.CreateIndex(index).BodyString(UserModel{}.Mapping()).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(createIndex.Index, "索引创建成功")
}
