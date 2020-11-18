package goInterventionImage

import (
	"fmt"
	"testing"
)

func TestInterventionImage_SaveToWEBP_SaveToWEBP(t *testing.T) {
	img, err := NewInterventionImage(&Config{
		FilePath: "./images/test.jpg",
	})
	if err != nil {
		panic(err)
	}
	img.SetWaterMark(&WaterMarkConfig{
		DestPosition: " left  buttom ",
	})
	if err := img.AddWaterMarkImg(""); err != nil {
		fmt.Println(err.Error())
	}
	img.SaveToWEBP("13212", 80)
}

func TestInterventionImage_Resize(t *testing.T) {
	img, err := NewInterventionImage(&Config{
		FilePath: "./images/test.jpg",
	})
	if err != nil {
		panic(err)
	}
	img.Resize(300, 600)
	if err := img.AddWaterMarkText("this is a testing"); err != nil {
		fmt.Println(err.Error())
	}
	img.SaveToPNG("123")
}

func TestInterventionImage_Save(t *testing.T) {
	img, err := NewInterventionImage(&Config{})
	if err != nil {
		panic(err)
	}
	if err := img.Resize(300, 600); err != nil {
		fmt.Println(err.Error())
	}
	if err := img.AddWaterMarkText("this is a testing"); err != nil {
		fmt.Println(err.Error())
	}
	img.SaveToJPG("123", 80)
}