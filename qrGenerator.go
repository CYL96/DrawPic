package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	qrcode "github.com/skip2/go-qrcode"
)

//  -------------------------------------- 创建二维码 ------------------------------------------------
/*
//  BuildQRCode
//  @Description:
//  @param [content]: 内容
//  @param [dir]: 保存目录
//  @param [name]: 保存文件名
//  @param [backCo]: 背景色
//  @param [forgeCo]: 前景色
//  @return [string]: 完整的文件路径
//  @return [error]: 错误信息
*/
func BuildQRCode(content, dir, name string, backCo, forgeCo color.Color) (string, error) {
	if len(dir) > 0 && len(dir) != strings.LastIndex(dir, "/")+1 {
		dir += "/"
	}
	file := dir + name
	// strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg";

	// 生成二维码;
	qr, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		return "", err
	}
	if backCo != nil {
		qr.BackgroundColor = backCo
	}
	if forgeCo != nil {
		qr.ForegroundColor = forgeCo
	}
	err = qr.WriteFile(256, file)
	if err != nil {
		return "", err
	}
	return file, nil
}

//  -------------------------------------- HEX 转 RBGA ------------------------------------------------
/*
//  HexToRGBA
//  @Description:  HEX 转 RBGA
//  @param [co]: HEX
//  @return [color.RGBA]: RBGA
*/
func HexToRGBA(co string) color.RGBA {
	if len(co) < 6 {
		return color.RGBA{}
	}
	r, _ := strconv.ParseUint(co[:2], 16, 8)
	g, _ := strconv.ParseUint(co[2:4], 16, 8)
	b, _ := strconv.ParseUint(co[4:], 16, 8)
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}

//  -------------------------------------- 生成二维码 ------------------------------------------------
/*
//  QrLogoCode
//  @Description:
//  @param [content]: 内容
//  @param [logo]: 底图logo
//  @return [image.Image]: 二维码
//  @return [error]: 错误信息
*/
func QrLogoCode(content string, logo string) (image.Image, error) {
	var (
		bgImg      image.Image
		offset     image.Point
		avatarFile *os.File
		avatarImg  image.Image
		err        error
	)

	bgImg, err = createQrCode(content)

	if err != nil {
		return nil, err
	}
	avatarFile, err = os.Open(logo)
	avatarImg, err = png.Decode(avatarFile)
	avatarImg = ImageResize(avatarImg, 40, 40)
	b := bgImg.Bounds()
	// 设置为居中
	offset = image.Pt((b.Max.X-avatarImg.Bounds().Max.X)/2, (b.Max.Y-avatarImg.Bounds().Max.Y)/2)
	m := image.NewRGBA(b)
	draw.Draw(m, b, bgImg, image.Point{X: 0, Y: 0}, draw.Src)
	draw.Draw(m, avatarImg.Bounds().Add(offset), avatarImg, image.Point{X: 0, Y: 0}, draw.Over)
	return m, err
}

//  -------------------------------------- 创建二维码 ------------------------------------------------
/*
//  createQrCode
//  @Description:
//  @param [content]: 内容
//  @return [img]: 二维码
//  @return [err]: 错误信息
*/
func createQrCode(content string) (img image.Image, err error) {
	var qrCode *qrcode.QRCode
	qrCode, err = qrcode.New(content, qrcode.Highest)
	if err != nil {
		return nil, err
	}
	qrCode.DisableBorder = true
	img = qrCode.Image(150)
	return img, nil
}

//  -------------------------------------- 二维码大小调整 ------------------------------------------------
/*
//  ImageResize
//  @Description: 二维码大小调整
//  @param [src]: 二维码图片
//  @param [w]: 宽度
//  @param [h]: 高度
//  @return [image.Image]: 调整后的二维码
*/
func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}
