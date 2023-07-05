// 验证码生成库
package comm

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/big"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/yuanzhi-ai/luban/server/log"
)

const (
	defaultWidth = 310
	defaultHight = 155
	VerifyLen    = 6
	chars        = "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	charsLen     = len(chars)
	fontPath     = "./app_data/captcha/fonts/"
	fontName     = "DENNEthree-dee.ttf"
	fontSize     = 12
)

// 图像验证码
// W图像宽度，H图像高度
type CaptchaGenerator struct {
	w       int
	h       int
	codeLen int
}

// GetNewCaptchaGenerator 获取一个验证码生成器
func GetNewCaptchaGenerator() *CaptchaGenerator {
	cg := CaptchaGenerator{w: defaultWidth, h: defaultHight, codeLen: VerifyLen}
	return &cg
}

// GetVerCode 获取一个验证码
// 返回一个验证码的id,和base64编码后的验证码图像
// 验证码的id = md5(验证码答案+秘钥)
// 故此解密只要再使用用户答案做相同计算与id比较即可
// 返回 验证码id， 验证码的url，和是否出错
func (cg *CaptchaGenerator) GetCaptcha() (string, string, error) {
	// 生成一个随机验证码
	captcha := cg.getRandStr()
	skeyInstance := GetSkeyInstance()
	capSkey, err := skeyInstance.GetSkey(CaptchaSkey)
	if err != nil || capSkey == "" {
		log.Errorf("get capSkey err:%v", err)
		return "", "", err
	}
	// 使用md5加密验证码和验证码的skey，作为改验证码的id
	capId := Md5Encode(captcha + capSkey)
	// 画图
	img := cg.initCanvas()
	capImg, err := cg.doImage(captcha, img)
	if err != nil {
		return "", "", err
	}
	// capImg 转base64 url
	encodeImg := cg.base64Img(capImg)
	return capId, encodeImg, nil
}

func (cg *CaptchaGenerator) VerifyCode(capId string, userAnswer string) (bool, error) {
	skeyInstance := GetSkeyInstance()
	capSkey, err := skeyInstance.GetSkey(CaptchaSkey)
	if err != nil || capSkey == "" {
		log.Errorf("get capSkey err:%v", err)
		return false, err
	}
	singAnswer := Md5Encode(userAnswer + capSkey)
	if capId != singAnswer {
		return false, nil
	}
	return true, nil
}

// 获取区间[-m, n]的随机数
func (cg *CaptchaGenerator) RangeRand(min, max int64) int64 {
	left := min
	right := max
	if min > max {
		left = max
		right = min
	}
	if left < 0 {
		f64Min := math.Abs(float64(left))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(right+1+i64Min))
		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(right-left+1))

		return left + result.Int64()
	}
}

// 初始化画布
func (cg *CaptchaGenerator) initCanvas() *image.RGBA {
	dest := image.NewRGBA(image.Rect(0, 0, cg.w, cg.h))

	// 随机色
	r := uint8(255) // uint8(captcha.RangeRand(50, 250))
	g := uint8(255) // uint8(captcha.RangeRand(50, 250))
	b := uint8(255) // uint8(captcha.RangeRand(50, 250))

	// 填充背景色
	for x := 0; x < cg.w; x++ {
		for y := 0; y < cg.h; y++ {
			dest.Set(x, y, color.RGBA{r, g, b, 255}) //设定alpha图片的透明度
		}
	}

	return dest
}

// getRandStr 获取一个验证码
func (cg *CaptchaGenerator) getRandStr() (randStr string) {

	for i := 0; i < cg.codeLen; i++ {
		randIndex := cg.RangeRand(0, int64(len(chars)-1))
		randStr += chars[randIndex : randIndex+1]
	}
	return randStr
}

// doImage 画图
func (cg *CaptchaGenerator) doImage(code string, dest *image.RGBA) (*image.RGBA, error) {
	gc := draw2dimg.NewGraphicContext(dest)
	defer gc.Close()
	defer gc.FillStroke()
	err := cg.setFont(gc)
	if err != nil {
		log.Errorf("set font err:%v", err)
		return nil, err
	}
	cg.doPoint(gc)
	cg.doLine(gc)
	cg.doSinLine(gc)
	cg.doCode(gc, code)
	return dest, nil
}

// 设置相关字体
func (cg *CaptchaGenerator) setFont(gc *draw2dimg.GraphicContext) error {

	// 字体文件
	fontFile := strings.TrimRight(fontPath, "/") + "/" + strings.TrimLeft(fontName, "/")

	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		log.Errorf("read font err:%v", err)
		return err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Errorf("parseFont err:%v", err)
		return err
	}

	// 设置自定义字体相关信息
	gc.FontCache = draw2d.NewSyncFolderFontCache(fontPath)
	gc.FontCache.Store(draw2d.FontData{Name: fontName, Family: 0, Style: draw2d.FontStyleNormal}, font)
	gc.SetFontData(draw2d.FontData{Name: fontName, Style: draw2d.FontStyleNormal})

	// 设置字体大小
	gc.SetFontSize(fontSize)
	return nil
}

// 增加干扰点
func (cg *CaptchaGenerator) doPoint(gc *draw2dimg.GraphicContext) {
	for n := 0; n < 50; n++ {
		gc.SetLineWidth(float64(cg.RangeRand(1, 3)))

		// 随机色
		r := uint8(cg.RangeRand(0, 255))
		g := uint8(cg.RangeRand(0, 255))
		b := uint8(cg.RangeRand(0, 255))

		gc.SetStrokeColor(color.RGBA{r, g, b, 255})

		x := cg.RangeRand(0, int64(cg.w)+10) + 1
		y := cg.RangeRand(0, int64(cg.h)+5) + 1

		gc.MoveTo(float64(x), float64(y))
		gc.LineTo(float64(x+cg.RangeRand(1, 2)), float64(y+cg.RangeRand(1, 2)))

		gc.Stroke()
	}
}

// 增加干扰线
func (cg *CaptchaGenerator) doLine(gc *draw2dimg.GraphicContext) {
	// 设置干扰线
	for n := 0; n < 5; n++ {
		// gc.SetLineWidth(float64(captcha.RangeRand(1, 2)))
		gc.SetLineWidth(1)

		// 随机背景色
		r := uint8(cg.RangeRand(0, 255))
		g := uint8(cg.RangeRand(0, 255))
		b := uint8(cg.RangeRand(0, 255))

		gc.SetStrokeColor(color.RGBA{r, g, b, 255})

		// 初始化位置
		gc.MoveTo(float64(cg.RangeRand(0, int64(cg.w)+10)), float64(cg.RangeRand(0, int64(cg.h)+5)))
		gc.LineTo(float64(cg.RangeRand(0, int64(cg.w)+10)), float64(cg.RangeRand(0, int64(cg.h)+5)))

		gc.Stroke()
	}
}

// 增加正弦干扰线
func (cg *CaptchaGenerator) doSinLine(gc *draw2dimg.GraphicContext) {
	h1 := cg.RangeRand(-12, 12)
	h2 := cg.RangeRand(-1, 1)
	w2 := cg.RangeRand(5, 20)
	h3 := cg.RangeRand(5, 10)

	h := float64(cg.h)
	w := float64(cg.w)

	// 随机色
	r := uint8(cg.RangeRand(128, 255))
	g := uint8(cg.RangeRand(128, 255))
	b := uint8(cg.RangeRand(128, 255))

	gc.SetStrokeColor(color.RGBA{r, g, b, 255})
	gc.SetLineWidth(float64(cg.RangeRand(2, 4)))

	var i float64
	for i = -w / 2; i < w/2; i = i + 0.1 {
		y := h/float64(h3)*math.Sin(i/float64(w2)) + h/2 + float64(h1)

		gc.LineTo(i+w/2, y)

		if h2 == 0 {
			gc.LineTo(i+w/2, y+float64(h2))
		}
	}

	gc.Stroke()
}

// 验证码字符设置到图像上
func (cg *CaptchaGenerator) doCode(gc *draw2dimg.GraphicContext, code string) {
	for l := 0; l < len(code); l++ {
		y := cg.RangeRand(int64(fontSize)-1, int64(cg.h)+6)
		x := cg.RangeRand(1, 20)

		// 随机色
		r := uint8(cg.RangeRand(0, 200))
		g := uint8(cg.RangeRand(0, 200))
		b := uint8(cg.RangeRand(0, 200))

		gc.SetFillColor(color.RGBA{r, g, b, 255})
		gc.FillStringAt(string(code[l]), float64(x)+fontSize*float64(l), float64(int64(cg.h)-y)+fontSize)
		gc.Stroke()
	}
}

// 将图片base64转码
func (cg *CaptchaGenerator) base64Img(dest *image.RGBA) string {

	emptyBuff := bytes.NewBuffer(nil)
	png.Encode(emptyBuff, dest)
	dist := make([]byte, 50000)
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())
	index := bytes.IndexByte(dist, 0)
	baseImage := string(dist[0:index])
	return baseImage
}
