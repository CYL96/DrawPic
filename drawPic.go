package main

import (
	"bytes"
	"encoding/base64"
	// "fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/bmp"
	"golang.org/x/image/math/fixed"
)

type DrawPicStruct struct {
	DefaultPic image.Image       // 默认图片
	Mgba       *image.NRGBA      // 画布
	Font       *freetype.Context // 字体
}

// 图片类型
type PicType int

const (
	PIC_TYPE_JPEG = PicType(0) // jpeg
	PIC_TYPE_PNG  = PicType(1) // png
	PIC_TYPE_BMP  = PicType(2) // bmp
)

//  -------------------------------------- 创建画布 ------------------------------------------------
/*
//  NewDrawPic
//  @Description: 通过图片路径创建画布
//  @param [defaultPic]: 图片路径
//  @return [d]: 画布
//  @return [err]: 错误信息
*/
func NewDrawPic(picPath string) (d *DrawPicStruct, err error) {
	var (
		imgb *os.File
		img  image.Image
	)
	imgb, err = os.Open(picPath)
	if err != nil {
		return
	}
	defer imgb.Close()
	img, _, err = image.Decode(imgb)
	if err != nil {
		return
	}
	d = new(DrawPicStruct)
	d.DefaultPic = img
	return
}

//  -------------------------------------- 通过base64创建画布 ------------------------------------------------
/*
//  NewDrawPicBase64
//  @Description:  通过base64创建画布
//  @param [bs64]: 图片的base64
//  @return [d]: 画布
//  @return [err]: 错误信息
*/
func NewDrawPicBase64(bs64 string) (d *DrawPicStruct, err error) {
	var (
		img  image.Image
		data []byte
	)
	data, err = base64.StdEncoding.DecodeString(bs64)
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	d = new(DrawPicStruct)
	d.DefaultPic = img
	return
}

//  -------------------------------------- 快捷创建二维码画布 ------------------------------------------------
/*
//  NewDrawQRCodePic
//  @Description: 生成二维码画布
//  @param [content]: 二维码内容
//  @param [size]: 二维码大小
//  @param [backCo]: 背景色
//  @param [forgeCo]: 前景色
//  @return [d]: 画布
//  @return [err]: 错误信息
*/
func NewDrawQRCodePic(content string, size int, backCo, forgeCo color.Color) (d *DrawPicStruct, err error) {
	var (
		qr *qrcode.QRCode
	)
	qr, err = qrcode.New(content, qrcode.Highest)
	if err != nil {
	}
	if backCo != nil {
		qr.BackgroundColor = backCo
	}
	if forgeCo != nil {
		qr.ForegroundColor = forgeCo
	}
	d = new(DrawPicStruct)
	d.DefaultPic = qr.Image(size)
	return
}

//  -------------------------------------- 图片大小调整 ------------------------------------------------
/*
//  Resize
//  @Description: 调整图片的大小。缩放
//  @receiver [this]: 画布
//  @param [width]: 宽度
//  @param [height]: 高度
*/
func (this *DrawPicStruct) Resize(width, height uint) {
	this.DefaultPic = resize.Resize(width, height, this.DefaultPic, resize.Lanczos3)
}

//  -------------------------------------- 初始化字体 ------------------------------------------------
/*
//  InitFront
//  @Description: 给画布创建字体
//  @receiver [this]: 画布
//  @param [font_path]: 字体路径
//  @return [err]: 错误信息
*/
func (this *DrawPicStruct) InitFront(font_path string) (err error) {
	var (
		fontBytes []byte
		font      *truetype.Font
	)
	// 打开字体
	fontBytes, err = ioutil.ReadFile(font_path)
	if err != nil {
		return
	}
	// 载入字体数据
	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		return
	}
	this.Font = freetype.NewContext()
	this.Font.SetFont(font)
	this.SetFontSize(26)
	this.SetFontDpi(72)
	this.SetFontColor(0, 0, 0, 255)
	return
}

//  -------------------------------------- 设置字体dpi ------------------------------------------------
/*
//  SetFontDpi
//  @Description: 设置字体dpi
//  @receiver [this]: 画布
//  @param [dpi]: dpi
*/
func (this *DrawPicStruct) SetFontDpi(dpi float64) {
	this.Font.SetDPI(dpi)
}

//  -------------------------------------- 设置字体的大小 ------------------------------------------------
/*
//  SetFontSize
//  @Description: 设置字体的大小
//  @receiver [this]: 画布
//  @param [size]: 字体大小
*/
func (this *DrawPicStruct) SetFontSize(size float64) {
	this.Font.SetFontSize(size)
}

//  -------------------------------------- 设置字体颜色 ------------------------------------------------
/*
//  SetFontColor
//  @Description: 具体RGBA的值查颜色对照表
//  @receiver [this]: 画布
//  @param [R]: R
//  @param [G]: G
//  @param [B]: B
//  @param [A]: A
*/
func (this *DrawPicStruct) SetFontColor(R, G, B, A uint8) {
	this.Font.SetSrc(image.NewUniform(color.RGBA{R, G, B, A}))
}

//  -------------------------------------- 在图上写字 ------------------------------------------------
/*
//  DrawFont
//  @Description: 在图上写字
//  @receiver [this]: 画布
//  @param [content]: 内容
//  @param [X]: 基于图片的X
//  @param [Y]: 基于图片的Y
//  @return [err]: 错误信息
*/
func (this *DrawPicStruct) DrawFont(content string, X, Y int) (err error) {
	var (
		offset fixed.Point26_6
		bounds image.Rectangle
		mgba   *image.NRGBA
	)
	offset = freetype.Pt(X, Y)
	bounds = this.DefaultPic.Bounds()
	mgba = image.NewNRGBA(bounds)
	draw.Draw(mgba, bounds, this.DefaultPic, image.ZP, draw.Src)
	this.Font.SetClip(bounds)
	this.Font.SetDst(mgba)

	_, err = this.Font.DrawString(content, offset)
	// fmt.Println(err)
	this.DefaultPic = mgba
	return
}

// 批量 在图上写字
type DrawFontBatchT struct {
	Content string
	X, Y    int
}

//  -------------------------------------- 基于图片的批量写字 ------------------------------------------------
/*
//  DrawFontBatch
//  @Description: 基于图片的批量写字
//  @receiver [this]: 画布
//  @param [ext]: 内容
//  @return [err]: 错误信息
*/
func (this *DrawPicStruct) DrawFontBatch(ext []DrawFontBatchT) (err error) {
	var (
		offset fixed.Point26_6
		bounds image.Rectangle
		mgba   *image.NRGBA
	)
	bounds = this.DefaultPic.Bounds()
	mgba = image.NewNRGBA(bounds)
	draw.Draw(mgba, bounds, this.DefaultPic, image.ZP, draw.Src)
	this.Font.SetClip(bounds)
	this.Font.SetDst(mgba)
	for _, t := range ext {
		offset = freetype.Pt(t.X, t.Y)
		_, err = this.Font.DrawString(t.Content, offset)
		if err != nil {
			return
		}
	}
	// fmt.Println(err)
	this.DefaultPic = mgba
	return
}

//  -------------------------------------- 画图 ------------------------------------------------
/*
//  DrawPic
//  @Description: 在图片上画图
//  @receiver [this]: 画布
//  @param [pic]: 图片
//  @param [X]: 基于图片的X
//  @param [Y]: 基于图片的Y
//  @return [err]: 错误信息
*/
func (this *DrawPicStruct) DrawPic(pic *DrawPicStruct, X, Y int) (err error) {
	var (
		offset image.Point
		bounds image.Rectangle
		mgba   *image.NRGBA
	)
	offset = image.Pt(X, Y)
	bounds = this.DefaultPic.Bounds()
	mgba = image.NewNRGBA(bounds)

	draw.Draw(mgba, bounds, this.DefaultPic, image.ZP, draw.Src)
	draw.Draw(mgba, pic.DefaultPic.Bounds().Add(offset), pic.DefaultPic, image.ZP, draw.Over)
	this.DefaultPic = mgba
	return
}

//  -------------------------------------- 保存 ------------------------------------------------
/*
//  Save
//  @Description: 保存图片
//  @receiver [this]: 画布
//  @param [path]: 保存路径
//  @param [quality]: 质量【0-100】
//  @param [picType]: 图片类型
//  @return [err]: 错误信息
*/
func (this *DrawPicStruct) Save(path string, quality int, picType PicType) (err error) {
	var (
		imgF *os.File
	)
	imgF, err = os.Create(path)
	if err != nil {
		return
	}
	defer imgF.Close()
	switch picType {
	case PIC_TYPE_JPEG:
		jpeg.Encode(imgF, this.DefaultPic, &jpeg.Options{Quality: quality})
	case PIC_TYPE_PNG:
		png.Encode(imgF, this.DefaultPic)
	case PIC_TYPE_BMP:
		bmp.Encode(imgF, this.DefaultPic)
	}
	return
}

//  -------------------------------------- 保存文件为base64 ------------------------------------------------
/*
//  SaveToBase64
//  @Description:
//  @receiver [this]: 画布
//  @param [quality]: 质量【0-100】
//  @param [picType]: 图片类型
//  @return [b64]: base64值
//  @return [err]: 错误信息
*/
func (this *DrawPicStruct) SaveToBase64(quality int, picType PicType) (b64 string, err error) {
	imgF := bytes.NewBuffer(nil)
	switch picType {
	case PIC_TYPE_JPEG:
		jpeg.Encode(imgF, this.DefaultPic, &jpeg.Options{Quality: quality})
	case PIC_TYPE_PNG:
		png.Encode(imgF, this.DefaultPic)
	case PIC_TYPE_BMP:
		bmp.Encode(imgF, this.DefaultPic)
	}

	b64 = base64.StdEncoding.EncodeToString(imgF.Bytes())
	return
}
