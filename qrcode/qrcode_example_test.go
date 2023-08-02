package qrcode_test

import (
	"fmt"
	"github.com/Is999/go-code-snippets/qrcode"
	"image/color"
	"net/http"
)

func ExampleQRCode() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 实例QRCode
		qr := qrcode.NewQRCode("https://github.com/Is999")

		// 设置相关参数
		qr.SetSize(72)                                                   // 设置图片大小
		qr.SetBackgroundColor(color.RGBA{R: 255, G: 215, B: 0, A: 255}). // 设置背景颜色
											SetForegroundColor(color.RGBA{R: 55, G: 245, B: 105, A: 255}) // 设置前景颜色
		qr.SetRecoveryLevel(2) // 设置容错级别

		// 创建二维码图片
		err := qr.CreateQrCode()
		if err != nil {
			fmt.Println("创建二维码识失败：", err.Error())
		}

		//// 给二维码填充色彩
		err = qr.FillColor("../images/color.png", 0, true)
		if err != nil {
			fmt.Println("给二维码填充色彩失败：", err.Error())
		}

		// 创建带logo的二维码
		err = qr.CreateAvatar("../images/logo.png", 20, 20, false)
		if err != nil {
			fmt.Println("创建带logo的二维码失败：", err.Error())
		}

		// 设置二维码边框
		err = qr.SetBorder(8, color.RGBA{R: 235, G: 255, B: 215, A: 255})
		if err != nil {
			fmt.Println("设置二维码边框失败：", err.Error())
		}

		//// 二维码与海报合成
		qr.SetHorizontal(11). // 设置偏移位置
					SetOffsetX(129) // 设置图片X轴方向偏移量
		err = qr.CreateAvatar("../images/new-year-poster.png", 338, 521, true)
		if err != nil {
			fmt.Println("二维码与海报合成失败：", err.Error())
		}

		// 显示图片
		err = qr.View(w)
		if err != nil {
			fmt.Println("显示图片失败：", err.Error())
		}
	})
}
