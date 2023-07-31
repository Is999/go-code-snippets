package qrcode

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
)

type QRCode struct {
	img              image.Image // 二维码图片
	size             int         // 二维码大小
	content          string      // 二维码内容
	backgroundColor  color.RGBA  // 背景颜色
	foregroundColor  color.RGBA  // 前景颜色
	recoveryLevel    uint8       // 容错级别 0-3：从低到高分别对应 Low 7%，Medium 15%，High 25%， Highest 30%
	horizontal       uint8       // 设置图片偏移量: 0自定义, 1居中, 2最左上, 3最右上, 4最右下, 5最左下, 6最上居中, 7最右居中, 8最下居中, 9最左居中, 10最右偏移, 11最下偏移
	offsetX, offsetY int         // 设置图片偏移量
}

// SetSize 二维码大小
func (q *QRCode) SetSize(size int) *QRCode {
	q.size = size
	return q
}

// SetBackgroundColor 背景颜色
func (q *QRCode) SetBackgroundColor(backgroundColor color.RGBA) *QRCode {
	q.backgroundColor = backgroundColor
	return q
}

// SetForegroundColor 前景颜色
func (q *QRCode) SetForegroundColor(foregroundColor color.RGBA) *QRCode {
	q.foregroundColor = foregroundColor
	return q
}

// SetRecoveryLevel 容错级别 0-3：从低到高分别对应 Low 7%，Medium 15%，High 25%， Highest 30%
func (q *QRCode) SetRecoveryLevel(recoveryLevel uint8) *QRCode {
	q.recoveryLevel = recoveryLevel
	return q
}

// SetHorizontal 设置图片偏移量（偏移方式）: 0自定义, 1居中, 2最左上, 3最右上, 4最右下, 5最左下, 6最上居中, 7最右居中, 8最下居中, 9最左居中, 10最右偏移, 11最下偏移
func (q *QRCode) SetHorizontal(horizontal uint8) *QRCode {
	q.horizontal = horizontal
	return q
}

// SetOffsetX 设置图片X轴方向偏移量：Horizontal为[0,11], 此设置生效, 否则不生效
func (q *QRCode) SetOffsetX(offsetX int) *QRCode {
	q.offsetX = offsetX
	return q
}

// SetOffsetY 设置图片Y轴方向偏移量：Horizontal为[0,10], 此设置生效, 否则不生效
func (q *QRCode) SetOffsetY(offsetY int) *QRCode {
	q.offsetY = offsetY
	return q
}

// SetContent 二维码内容
func (q *QRCode) SetContent(content string) *QRCode {
	q.content = content
	return q
}

func NewQRCode(content string) *QRCode {
	return &QRCode{
		img:             nil,
		size:            256, // 二维码图片大小：默认256px
		content:         content,
		backgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 背景色: 白色
		foregroundColor: color.RGBA{A: 255},                         // 前景色: 黑色
		recoveryLevel:   3,                                          // 容错等级：默认最大容错
		horizontal:      1,                                          // 偏移量居中
		offsetX:         0,                                          // Horizontal为[1-10], 此设置不生效
		offsetY:         0,                                          // Horizontal为[1-9,11], 此设置不生效
	}
}

// View 展示二维码图片
func (q *QRCode) View(writer io.Writer) error {
	// 判断是否已经创建了二维码， 没有则创建
	if q.img == nil {
		if err := q.CreateQrCode(); err != nil {
			return err
		}
	}

	return png.Encode(writer, q.img)
}

// CreateQrCode 创建二维码
func (q *QRCode) CreateQrCode() error {
	qrCode, err := qrcode.New(q.content, qrcode.RecoveryLevel(q.recoveryLevel))
	if err != nil {
		return fmt.Errorf("create qrcode err: %s" + err.Error())
	}

	qrCode.DisableBorder = true                // 禁用边框（默认不支持设置边框大小，可以改用 Border 设置边框）
	qrCode.BackgroundColor = q.backgroundColor // 背景色
	qrCode.ForegroundColor = q.foregroundColor // 前景色
	q.img = qrCode.Image(q.size)               // 设置大小

	return nil
}

// SetBorder 给二维码添加边框
// size 自定义边框大小
// color 背景色: 白色 color.RGBA{R: 255, G: 255, B: 255, A: 255}
func (q *QRCode) SetBorder(size uint8, colors ...color.Color) error {
	// 判断是否已经创建了二维码， 没有则创建
	if q.img == nil {
		if err := q.CreateQrCode(); err != nil {
			return err
		}
	}

	var col color.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	if len(colors) > 0 {
		col = colors[0]
	}

	q.mergeImg(q.img, q.createImage(int(size)+q.size, col))
	return nil
}

// CreateAvatar 创建一个带LOGO的二维码
func (q *QRCode) CreateAvatar(imgPath string, width, height uint, isBgImg bool) error {
	// 判断是否已经创建了二维码， 没有则创建
	if q.img == nil {
		if err := q.CreateQrCode(); err != nil {
			return err
		}
	}

	// 加载 AvatarImgPath
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("imgPath open: %s" + err.Error())
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("imgFile decode: %s " + err.Error())
	}

	var fgImg, bgImg image.Image
	//设置图片大小
	if isBgImg {
		fgImg, bgImg = q.img, img
		bgImg = resize.Resize(width, height, bgImg, resize.Lanczos3)
	} else {
		fgImg, bgImg = img, q.img
		fgImg = resize.Resize(width, height, fgImg, resize.Lanczos3)
	}

	q.mergeImg(fgImg, bgImg)
	return nil
}

// mergeImg 合并图片
func (q *QRCode) mergeImg(fgImg, bgImg image.Image) {
	fg := fgImg.Bounds()
	bg := bgImg.Bounds()

	m := image.NewRGBA(bg)
	draw.Draw(m, bg, bgImg, image.Point{X: 0, Y: 0}, draw.Src)
	draw.Draw(m, fg.Add(q.getImgPoint(fg, bg)), fgImg, image.Point{X: 0, Y: 0}, draw.Over)
	q.img = m
}

// getImgPoint 获取两个图片间的偏移量
func (q *QRCode) getImgPoint(a, b image.Rectangle) image.Point {
	switch q.horizontal {
	case 1:
		//1居中
		q.offsetX = (b.Max.X - a.Max.X) / 2
		q.offsetY = (b.Max.Y - a.Max.Y) / 2
	case 2:
		//2最左上
		q.offsetX = 0
		q.offsetY = 0
	case 3:
		//3最右上
		q.offsetX = b.Max.X - a.Max.X
		q.offsetY = 0
	case 4:
		//4最右下
		q.offsetX = b.Max.X - a.Max.X
		q.offsetY = b.Max.Y - a.Max.Y
	case 5:
		//5最左下
		q.offsetX = 0
		q.offsetY = b.Max.Y - a.Max.Y
	case 6:
		//6最上居中
		q.offsetX = (b.Max.X - a.Max.X) / 2
		q.offsetY = 0
	case 7:
		//7最右居中
		q.offsetX = b.Max.X - a.Max.X
		q.offsetY = (b.Max.Y - a.Max.Y) / 2
	case 8:
		//8最下居中
		q.offsetX = (b.Max.X - a.Max.X) / 2
		q.offsetY = b.Max.Y - a.Max.Y
	case 9:
		//9最左居中
		q.offsetX = 0
		q.offsetY = (b.Max.Y - a.Max.Y) / 2
	case 10:
		//10最右偏移
		q.offsetX = b.Max.X - a.Max.X
	case 11:
		//11最下偏移
		q.offsetY = b.Max.Y - a.Max.Y
	}
	return image.Pt(q.offsetX, q.offsetY)
}

// 创建图片
func (q *QRCode) createImage(imgSize int, imgColor color.Color) image.Image {
	rect := image.Rectangle{Min: image.Point{}, Max: image.Point{X: imgSize, Y: imgSize}}
	p := color.Palette([]color.Color{imgColor, imgColor})
	return image.NewPaletted(rect, p)
}

// FillColor 填充二维码颜色
func (q *QRCode) FillColor(imgPath string, borderSize uint8, isFillBg bool) error {
	switch q.img.(type) {
	case *image.NRGBA:
	case *image.RGBA:
	default:
		//容错处理
		fg := q.img.Bounds()
		m := image.NewRGBA(fg)
		draw.Draw(m, fg, q.img, image.Point{X: 0, Y: 0}, draw.Src)
		q.img = m
	}

	// 获取二维码的宽高
	width, height := q.img.Bounds().Max.X, q.img.Bounds().Max.Y

	// 打开要填充的图片
	bgFile, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("fill color image open err: %s", err.Error())
	}

	defer bgFile.Close()

	// 将填充图解码成png图片
	bgImg, _, err := image.Decode(bgFile)
	if err != nil {
		return fmt.Errorf("fill color image decode err: %s", err.Error())
	}

	// 获取填充图的宽高
	bgWidth, bgHeight := bgImg.Bounds().Max.X, bgImg.Bounds().Max.Y

	// 检测二维码和填充图宽高是否一致
	if width != bgWidth || height != bgHeight {
		// 如果不一致将填充图剪裁
		bgImg = resize.Resize(uint(width), uint(height), bgImg, resize.Lanczos3)
	}
	imgClolor := q.foregroundColor
	if isFillBg {
		imgClolor = q.backgroundColor
	}
	// 开始填充二维码
	size := int(borderSize)
	imgSize := q.size - size*2
	for y := 0; y < q.img.Bounds().Max.X; y++ {
		for x := 0; x < q.img.Bounds().Max.X; x++ {
			if size > 0 && (x < size || x > size+imgSize || y < size || y > size+imgSize) {
				continue
			}
			qrImgColor := q.img.At(x, y)

			// 检测图片颜色 如果rgb值是 255 255 255 255 则像素点为白色 跳过
			// 如果rgba值是 0 0 0 0 则为透明色 跳过
			switch q.img.(type) {
			case *image.NRGBA:
				c := qrImgColor.(color.NRGBA)
				if (c.R == 0 && c.G == 0 && c.B == 0 && c.A == 0) || (c.R == imgClolor.R && c.G == imgClolor.G && c.B == imgClolor.B && c.A == imgClolor.A) {
					continue
				}
			case *image.RGBA:
				c := qrImgColor.(color.RGBA)
				if (c.R == 0 && c.G == 0 && c.B == 0 && c.A == 0) || (c.R == imgClolor.R && c.G == imgClolor.G && c.B == imgClolor.B && c.A == imgClolor.A) {
					continue
				}
			}

			// 获取要填充的图片的颜色
			bgImgColor := bgImg.At(x, y)

			// 填充颜色
			switch bgImg.(type) {
			case *image.RGBA64:
				c := bgImgColor.(color.RGBA64)
				q.img.(draw.Image).Set(x, y, color.RGBA64{R: c.R, G: c.G, B: c.B, A: c.A})

			case *image.NRGBA:
				c := bgImgColor.(color.NRGBA)
				q.img.(draw.Image).Set(x, y, color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A})

			case *image.RGBA:
				c := bgImgColor.(color.RGBA)
				q.img.(draw.Image).Set(x, y, color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A})

			case *image.YCbCr:
				c := bgImgColor.(color.YCbCr)
				q.img.(draw.Image).Set(x, y, color.YCbCr{Y: c.Y, Cb: c.Cb, Cr: c.Cr})
			default:
				return errors.New("no matching color format")
			}
		}
	}
	return nil
}
