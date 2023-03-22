package main

import (
	"fmt"
	"testing"
)

func TestNewDrawPic(t *testing.T) {
	// 初始化图片
	pic1, err := NewDrawPic("1.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 加在地址
	err = pic1.InitFront("./font/wryh.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 初始化图片
	pic2, err := NewDrawPic("2.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置图片的大小
	pic2.Resize(300, 300)
	// 画图
	pic1.DrawPic(pic2, 0, 100)
	pic2.Resize(300, 300)
	pic2.Resize(150, 150)
	pic1.DrawPic(pic2, 400, 200)

	// 写字
	pic1.DrawFont("hahahah", 300, 200)
	// 设置字体属性
	pic1.SetFontColor(255, 0, 0, 255)
	pic1.SetFontSize(40)
	pic1.DrawFont("hahahah", 400, 400)
	// 保存
	pic1.Save("3.png", 100, PIC_TYPE_PNG)
}
