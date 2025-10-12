package main

import (
	"fmt"
	"github.com/olivere/elastic/v7"
)

var (
	ESClient *elastic.Client
)

func EsConnect() {

	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("", ""),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connect success: ", client)
	ESClient = client
}

func main() {
	EsConnect()

	//fmt.Println(global.ESClient)

	CreateIndex()

	//docs.DocDelete()
	//docs.DocDeleteBatch()
	DocCreateBatch()
	DocFind()
	//docs.DocUpdate()

	select {}
}
