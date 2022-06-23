package main

import (
	"fmt"
	"image"
	"time"
)

type PositionConfig struct {
	Latitude  float64
	Longitude float64
	MapID     int
}

type NowcastComposer struct {
	Pos PositionConfig
}

func NewNowcastComposer(pos PositionConfig) NowcastComposer {
	return NowcastComposer{Pos: pos}
}

type NowcastImage struct {
	Image image.Image
	Date  time.Time
	Url   string
}

type NowcastImages []NowcastImage

func CreateGIF(pos PositionConfig) (image.Image, error) {
	composer := NewNowcastComposer(pos)
	return composer.downloadGIF()
}

func (c *NowcastComposer) downloadGIF() (image.Image, error) {
	now := time.Now()
	needDownload, err := DownloadableImage(now)
	if err != nil {
		return nil, fmt.Errorf("downloadable images : %w", err)
	}
	images, err := needDownload.Download()
	if err != nil {
		return nil, fmt.Errorf("download images : %w", err)
	}
	return images.composeGIF()
}

func (images NowcastImages) composeGIF() (image.Image, error) {
	//TODO
	return nil, nil
}

func (d NowcastImages) Download() (NowcastImages, error) {
	//TODO
	return NowcastImages{}, nil
}

func DownloadableImage(now time.Time, mapID int) (NowcastImages, error) {
	//TODO
	date := now.Format("20060102150400")
	nearMax := 12
	for i := 0; i < nearMax; i++ {
		fmt.Sprintf("https://www.jma.go.jp/bosai/rain/data/ra/20220620115000/rain01_%s_f%02d_a%02d.png", date, i, mapID)
	}
	return NowcastImages{}, nil
}

func DownloadImage(url string) (image.Image, error) {
	//TOOD
	return nil, nil
}

func MergeImage(img, on string) (image.Image, error) {
	//TODO
	return nil, nil
}
