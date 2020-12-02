# Go-Intervention-Image

GO Intervention Image is a golang image handling and manipulation library providing an easier and expressive way to create, edit, and compose images. 

[![Gitter](https://badges.gitter.im/Go-Intervention-Image/community.svg)](https://gitter.im/Go-Intervention-Image/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## Installation
```golang
go get github.com/xiyoudaodao/go-intervention-image
```

## Code Examples

```golang
// open an image file
img, err := NewInterventionImage(&Config{
    FilePath: "./images/test.jpg",
})
if err != nil {
    panic(err)
}

// set watermark
img.AddWaterMarkImg("", &WaterMarkConfig{
    DestPosition: " left  top ",
})
img.AddWaterMarkImg("", &WaterMarkConfig{
    DestPosition: " right  buttom ",
})
img.AddWaterMarkText("this is a testing", &WaterMarkConfig{
    DestPosition: " right  top ",
})
img.AddWaterMarkText("this is a testing", &WaterMarkConfig{
    DestPosition: " left  buttom ",
})

// resize image instance
img.Resize(300, 600)

// save image in desired format, save default jpeg
img.Save("test", 80)
img.SaveToJPG("test", 80)
img.SaveToPNG("test")

//or save to stream, default webp
img.SaveToStream(80)

//creates valid code
img, err := NewInterventionImage(nil)
if err != nil {
    panic(err)
}
verificationCode, _ := img.MakeVerificationCode(6, 100, 30)
imageByte, _ := img.SaveToWEBPStream(80)
//Display images in HTML
<img src=`data:image/webp;base64,${imageByte}`>
```

## Configuration
```golang
//Used to initialize the configuration
//There are (Image/NewNRGBA/FilePath) three ways to initialize a picture. If not, a blank picture with a black background is created.
//SaveFilefolder used to set the folder for file saving. The default is the program execution directory.
type Config struct {
	Image          image.Image
	NewNRGBA       *image.NRGBA
	FilePath       string
	SaveFilefolder string
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
```

## Contributing

Contributions to the Go Intervention Image library are welcome. 

## License

© xiyoudaodao, 2020~time.Now

Released under the [MIT License](https://github.com/xiyoudaodao/go-intervention-image/blob/master/License)
