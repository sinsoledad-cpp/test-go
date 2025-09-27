package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"

	"github.com/wenlng/go-captcha/v2/base/codec"
	"github.com/wenlng/go-captcha/v2/rotate"
)

var rotateCapt rotate.Captcha

func init() {
	builder := rotate.NewBuilder()

	// 可以使用预置的素材资源：https://github.com/wenlng/go-captcha-assets
	bgImage, err := loadPng("../resources/bg.png")
	if err != nil {
		log.Fatalln(err)
	}

	bgImage1, err := loadPng("../resources/bg1.png")
	if err != nil {
		log.Fatalln(err)
	}

	builder.SetResources(
		rotate.WithImages([]image.Image{
			bgImage,
			bgImage1,
		}),
	)

	rotateCapt = builder.Make()
}

func main() {
	captData, err := rotateCapt.Generate()
	if err != nil {
		log.Fatalln(err)
	}

	blockData := captData.GetData()
	if blockData == nil {
		log.Fatalln(">>>>> generate err")
	}

	block, _ := json.Marshal(blockData)
	fmt.Println(">>>>>", string(block))

	var mBase64, tBase64 string
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(">>>>> ", mBase64)
	fmt.Println(">>>>> ", tBase64)

	//err = captData.GetMasterImage().SaveToFile("../resources/master.png")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err = captData.GetThumbImage().SaveToFile("../resources/thumb.png")
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func loadPng(p string) (image.Image, error) {
	imgBytes, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return codec.DecodeByteToPng(imgBytes)
}
