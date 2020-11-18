package goInterventionImage

import (
	"testing"
)

func TestInterventionImage_SaveToWEBP_SaveToWEBP(t *testing.T) {
	img, err := NewInterventionImage(&Config{
		FilePath: "./images/test.jpg",
	})
	if err != nil {
		panic(err)
	}
	img.AddWaterMarkImg("", &WaterMarkConfig{
		DestPosition: " left  top ",
	})
	img.AddWaterMarkImg("", &WaterMarkConfig{
		DestPosition: " right  buttom ",
	})
	img.AddWaterMarkText("这是一个测试", &WaterMarkConfig{
		DestPosition: " right  top ",
	})
	img.AddWaterMarkText("这是一个测试", &WaterMarkConfig{
		DestPosition: " left  buttom ",
	})
	img.Resize(300, 600)
	img.SaveToJPG("13212", 80)
}

func TestInterventionImage_Resize(t *testing.T) {
	img, err := NewInterventionImage(&Config{
		FilePath: "./images/test.jpg",
	})
	if err != nil {
		panic(err)
	}
	img.AddWaterMarkText("this is a testing", nil)
	img.SaveToPNG("123")
}

func TestInterventionImage_Save(t *testing.T) {
	img, err := NewInterventionImage(nil)
	if err != nil {
		panic(err)
	}
	img.Resize(300, 600)
	img.AddWaterMarkText("this is a testing", nil)
	img.SaveToJPG("123", 80)
}
