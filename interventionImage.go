package goInterventionImage

import (
	"bytes"
	"fmt"
	webp2 "github.com/chai2010/webp"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type InterventionImage struct {
	fileByte        *bytes.Buffer
	image           image.Image
	newNRGBA        *image.NRGBA
	font            *truetype.Font
	filePath        string
	saveFileFolder  string
	waterMarkConfig *WaterMarkConfig
}

//打开图片
func (i *InterventionImage) openImageByte() (err error) {
	filetype, err := ioutil.ReadFile(i.filePath)
	if err != nil {
		log.Println(err)
		return
	}
	i.fileByte = bytes.NewBuffer(filetype)
	return
}

func (i *InterventionImage) setImage() (err error) {
	i.image, _, err = image.Decode(i.fileByte)
	return
}

//图片缩放
func (i *InterventionImage) Resize(dstW uint, dstH uint) (err error) {
	i.image = resize.Resize(dstW, dstH, i.newNRGBA, resize.Lanczos3)
	i.initNewNRGBA()
	return nil
}

//初始化
func (i *InterventionImage) initNewNRGBA() {
	i.newNRGBA = image.NewNRGBA(i.image.Bounds())
	for y := 0; y < i.newNRGBA.Bounds().Dy(); y++ {
		for x := 0; x < i.newNRGBA.Bounds().Dx(); x++ {
			i.newNRGBA.Set(x, y, i.image.At(x, y))
		}
	}
}

//计算图片区域色块（黑或白）
func (i *InterventionImage) calculateImgColor(dx int, dy int, w int, h int) color.RGBA {
	var count float64
	var bright float64
	count = 0
	bright = 0
	for x := 0; x < i.newNRGBA.Bounds().Dx(); x++ {
		if x < dx || x > (dx+w) {
			continue
		}
		for y := 0; y < i.newNRGBA.Bounds().Dy(); y++ {
			if y > dy || y < (dy-h) {
				continue
			}
			r, g, b, _ := i.newNRGBA.At(x, y).RGBA()
			floatR, _ := strconv.ParseFloat(fmt.Sprint(r>>8), 64)
			floatG, _ := strconv.ParseFloat(fmt.Sprint(g>>8), 64)
			floatB, _ := strconv.ParseFloat(fmt.Sprint(b>>8), 64)
			count++
			bright = bright + 0.299*floatR + 0.587*floatG + 0.114*floatB
		}
	}
	if bright/count < 151 {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	return color.RGBA{R: 0, G: 0, B: 0, A: 255}
}

//Image watermark configuration item
//FontPath Custom font path, The format is TTF
//Fontbase64 Custom font path, TTF transfer base64
//FontSize Custom font size, The default value is 12
//FontDPI The default value is 72
//DestX, DestY Specifies where the watermark appears in the image
//DestPosition Describes where the watermark appears in the image，support('LEFT TOP'|'LEFT BUTTOM'|'RIGHT TOP'|'RIGHT BUTTOM'|'CENTER CENTER'),Case insensitive
//WaterMarkColor Specifies the watermark text color, If not specified, the text color (white and black) is set according to the background color of the picture
type WaterMarkConfig struct {
	FontPath       string
	Fontbase64     *[]byte
	FontSize       float64
	FontDPI        float64
	DestX          int
	DestY          int
	DestPosition   string
	WaterMarkColor color.RGBA
}

//设置水印
func (i *InterventionImage) SetWaterMark(w *WaterMarkConfig) (err error) {
	i.waterMarkConfig = w
	var fontBytes *[]byte
	if w.Fontbase64 != nil && len(*w.Fontbase64) > 0 {
		fontBytes = w.Fontbase64
	}
	if w.FontPath != "" {
		if *fontBytes, err = ioutil.ReadFile(w.FontPath); err != nil {
			return
		}
	}
	if fontBytes != nil && len(*fontBytes) > 0 {
		if i.font, err = freetype.ParseFont(*fontBytes); err != nil {
			return
		}
	}
	return
}

func (i *InterventionImage) calculateXY(waterMarkW int, waterMarkH int, imgW int, imgH int) (destX int, destY int) {
	destX = 0
	destY = waterMarkH
	log.Println(strings.ToUpper(DeleteExtraSpace(i.waterMarkConfig.DestPosition)))
	switch strings.ToUpper(DeleteExtraSpace(i.waterMarkConfig.DestPosition)) {
	case "LEFT TOP":
	case "LEFT BUTTOM":
		destY = imgH - waterMarkH
	case "RIGHT TOP":
		destX = imgW - waterMarkW
	case "RIGHT BUTTOM":
		destX = imgW - waterMarkW
		destY = imgH - waterMarkH
	case "CENTER CENTER":
		destX = (imgW - waterMarkW) / 2
		destY = (imgH - waterMarkH) / 2
	}
	log.Printf("destX: %d  destY: %d  \r\n", destX, destY)
	return
}

//图片添加文字水印
func (i *InterventionImage) AddWaterMarkText(waterMarkText string) (err error) {
	c := freetype.NewContext()
	c.SetFont(i.font)
	if i.waterMarkConfig.FontDPI == 0 {
		c.SetDPI(72)
	} else {
		c.SetDPI(i.waterMarkConfig.FontDPI)
	}
	if i.waterMarkConfig.FontSize == 0 {
		i.waterMarkConfig.FontSize = 12
	}
	c.SetFontSize(i.waterMarkConfig.FontSize)
	c.SetClip(i.newNRGBA.Bounds())
	c.SetDst(i.newNRGBA)

	dsX := 10
	dsY := i.newNRGBA.Bounds().Dy() - int(c.PointToFixed(i.waterMarkConfig.FontSize)>>6) - 10
	if i.waterMarkConfig.DestX > 0 && i.waterMarkConfig.DestY > 0 {
		dsX = i.waterMarkConfig.DestX
		dsY = i.waterMarkConfig.DestY
	}

	fw := int(c.PointToFixed(i.waterMarkConfig.FontSize)>>6) * len([]rune(waterMarkText))
	fh := int(c.PointToFixed(i.waterMarkConfig.FontSize) >> 6)
	if i.waterMarkConfig.DestPosition != "" {
		dsX, dsY = i.calculateXY(
			fw,
			fh,
			i.newNRGBA.Bounds().Dx(),
			i.newNRGBA.Bounds().Dy(),
		)
	}

	if (i.waterMarkConfig.WaterMarkColor == color.RGBA{}) {
		c.SetSrc(image.NewUniform(i.calculateImgColor(dsX, dsY, fw, fh)))
	} else {
		c.SetSrc(image.NewUniform(i.waterMarkConfig.WaterMarkColor))
	}
	pt := freetype.Pt(dsX, dsY)
	_, err = c.DrawString(waterMarkText, pt)
	return err
}

//图片添加图片水印
func (i *InterventionImage) AddWaterMarkImg(imagePath string) (err error) {
	var waterMarkFile = bytes.NewBuffer(imgBase64)
	if imagePath != "" {
		wf, err := ioutil.ReadFile(imagePath)
		if err != nil {
			log.Println("AddWaterMarkImg open pngfile err.")
			return err
		}
		waterMarkFile = bytes.NewBuffer(wf)
	}
	waterMarkImage, _, err := image.Decode(waterMarkFile)
	if err != nil {
		log.Println("AddWaterMarkImg Decode image err.")
		return err
	}

	dsX, dsY := i.calculateXY(
		waterMarkImage.Bounds().Dx(),
		waterMarkImage.Bounds().Dy(),
		i.newNRGBA.Bounds().Dx(),
		i.newNRGBA.Bounds().Dy(),
	)
	offset := image.Pt(dsX, dsY)
	b := i.newNRGBA.Bounds()

	//image.ZP代表Point结构体，目标的源点，即(0,0)
	//draw.Src源图像透过遮罩后，替换掉目标图像
	//draw.Over源图像透过遮罩后，覆盖在目标图像上（类似图层）
	draw.Draw(i.newNRGBA, b, i.newNRGBA, image.ZP, draw.Src)
	draw.Draw(i.newNRGBA, waterMarkImage.Bounds().Add(offset), waterMarkImage, image.ZP, draw.Over)
	return err
}

func (i *InterventionImage) getFileName(filename string) (name string) {
	if filename == "" {
		filename = i.filePath
	}
	//去掉名称中的后缀名
	filename = strings.Replace(filename, filepath.Ext(filename), "", -1)
	name = filename
	return
}

func (i *InterventionImage) SaveToBMP(filename string) (string, error) {
	var path = i.saveFileFolder + i.getFileName(filename) + ".bmp"
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return path, err
	}
	err = bmp.Encode(out, i.newNRGBA)
	return path, err
}

func (i *InterventionImage) SaveToGIF(filename string) (string, error) {
	var path = i.saveFileFolder + i.getFileName(filename) + ".gif"
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return path, err
	}
	err = gif.Encode(out, i.newNRGBA, &gif.Options{})
	return path, err
}

func (i *InterventionImage) SaveToPNG(filename string) (string, error) {
	var path = i.saveFileFolder + i.getFileName(filename) + ".png"
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return path, err
	}
	err = png.Encode(out, i.newNRGBA)
	return path, err
}

func (i *InterventionImage) SaveToJPG(filename string, quality int) (string, error) {
	var path = i.saveFileFolder + i.getFileName(filename) + ".jpg"
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return path, err
	}
	err = jpeg.Encode(out, i.newNRGBA, &jpeg.Options{quality})
	return path, err
}

func (i *InterventionImage) SaveToWEBP(filename string, quality float32) (string, error) {
	var path = i.saveFileFolder + i.getFileName(filename) + ".webp"
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return path, err
	}
	err = webp2.Encode(out, i.newNRGBA, &webp2.Options{Lossless: false, Quality: quality})
	return path, err
}

//默认为jpg格式
func (i *InterventionImage) Save(filename string, quality int) (string, error) {
	return i.SaveToJPG(filename, quality)
}

func (i *InterventionImage) SaveToBMPStream() ([]byte, error) {
	var buf bytes.Buffer
	err := bmp.Encode(&buf, i.newNRGBA)
	return buf.Bytes(), err
}

func (i *InterventionImage) SaveToGIFStream() ([]byte, error) {
	var buf bytes.Buffer
	err := gif.Encode(&buf, i.newNRGBA, &gif.Options{})
	return buf.Bytes(), err
}

func (i *InterventionImage) SaveToJPGStream(quality int) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, i.newNRGBA, &jpeg.Options{Quality: quality})
	return buf.Bytes(), err
}

func (i *InterventionImage) SaveToPNGStream() ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, i.newNRGBA)
	return buf.Bytes(), err
}

func (i *InterventionImage) SaveToWEBPStream(quality float32) ([]byte, error) {
	var buf bytes.Buffer
	err := webp2.Encode(&buf, i.newNRGBA, &webp2.Options{Lossless: false, Quality: quality})
	return buf.Bytes(), err
}

//default is webp
func (i *InterventionImage) SaveToStream(quality float32) ([]byte, error) {
	return i.SaveToWEBPStream(quality)
}

//Used to initialize the configuration
//There are (Image/NewNRGBA/FilePath) three ways to initialize a picture. If not, a blank picture with a black background is created.
//SaveFilefolder used to set the folder for file saving. The default is the program execution directory.
type Config struct {
	Image          image.Image
	NewNRGBA       *image.NRGBA
	FilePath       string
	SaveFilefolder string
}

func NewInterventionImage(config *Config) (i *InterventionImage, err error) {
	i = &InterventionImage{
		waterMarkConfig: &WaterMarkConfig{},
	}
	if i.font, err = freetype.ParseFont(imgFontBase64); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if (*config == Config{}) {
		log.Println("no config, create blank image")
		var w, h = 800, 600
		i.newNRGBA = image.NewNRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				i.newNRGBA.Set(x, y, color.NRGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				})
			}
		}
	} else {
		if config.Image != nil {
			i.initNewNRGBA()
		}
		if config.NewNRGBA != nil {

		}
		if config.SaveFilefolder != "" {
			i.saveFileFolder = config.SaveFilefolder
			err = os.MkdirAll(config.SaveFilefolder, os.ModePerm)
			if err != nil {
				panic("MkdirAll fail.")
				return nil, err
			}
		} else {
			config.SaveFilefolder = "./"
		}
		if config.FilePath != "" {
			i.filePath = config.FilePath
			if i.openImageByte() != nil {
				panic("not found image file.")
				return nil, err
			}
			if i.setImage() != nil {
				panic("not found image type.")
				return nil, err
			}
			i.initNewNRGBA()
		}
	}
	return
}
