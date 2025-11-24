package stage7

import (
	stage4 "test-go/eino/stage04"
	stage5 "test-go/eino/stage05"
	stage6 "test-go/eino/stage06"
)

func BuildRAG() {
	docs := stage6.TransDoc()
	stage4.IndexerRAG(docs)
	results := stage5.RetrieverRAG("欲渡黄河冰塞川")
	for _, doc := range results {
		println(doc.ID)
		println("================================================")
		println(doc.Content)
	}
}
