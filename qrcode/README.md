<div align="center">
<h1>二维码</h1>
</div>

## 简介：

1. 基于go-qrcode的二次包装依赖以下三方包
   1. github.com/skip2/go-qrcode
   2. github.com/nfnt/resize
2. 支持二维码图片设置大小
3. 支持二维码嵌入logo
4. 支持二维码设置颜色（以图片为）
5. 支持二维码嵌入海报
6. 支持二维码相对logo或海报设置偏移位置：0 自定义, 1 居中, 2 最左上, 3 最右上, 4 最右下, 5 最左下, 6 最上居中, 7 最右居中, 8 最下居中, 9 最左居中, 10 最右偏移, 11 最下偏移

## 示例：

*[查看示例代码](https://github.com/Is999/go-code-snippets/blob/main/qrcode/qrcode_example_test.go#L10)*

<p align="center">
<img src="https://raw.githubusercontent.com/Is999/go-code-snippets/main/images/qrcode-default.png" alt="二维码">
<img src="https://raw.githubusercontent.com/Is999/go-code-snippets/main/images/qrcode-color.png" alt="彩色二维码">
<img src="https://raw.githubusercontent.com/Is999/go-code-snippets/main/images/qrcode-logo.png" alt="logo二维码">
<img src="https://raw.githubusercontent.com/Is999/go-code-snippets/main/images/qrcode-haibao.png" alt="海报二维码">
</p>
