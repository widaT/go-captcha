package main

//调整算法后，快速验证的工具
import (
	"encoding/base64"
	"io/ioutil"

	captcha "github.com/widaT/go-captcha/puzzle_captcha"
)

func main() {
	captcha.LoadBackgroudImages("./images/puzzle_captcha/backgroud")
	captcha.LoadBlockImages("./images/puzzle_captcha/block")

	ret, err := captcha.Run()
	if err != nil {
		return
	}

	saveImage("test_data/bg.png", ret.BackgroudImg)
	saveImage("test_data/bk.png", ret.BlockImg)
}

func saveImage(path string, base64Img string) error {
	i, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, i, 0666)
}
