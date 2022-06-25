package main

import (
	"fmt"
	"image/color"
	"os"
	"time"
)

type RainInfo struct {
	Rains []RainData
}

type RainData struct {
	Time      time.Time
	ImageID   int
	RainFall  int
	RainColor color.Color
}

var RainFallColor = map[string]int{
	"#f2f2ff": 1,
	"#a0d2ff": 5,
	"#218cff": 10,
	"#0041ff": 20,
	"#fff500": 30,
	"#ff9900": 50,
	"#ff2800": 80,
	"#b40068": 150,
}

func (n NowcastImage) SeekRainData(pos PositionConfig) (RainData, error) {
	rainFall := 0
	var rainColor color.Color
	for dy := 0; dy < pos.CursorHeight; dy++ {
		c := n.Image.At(pos.X, pos.Y+dy)
		r, g, b, a := c.RGBA()
		if a > 0 {
			colorCode := fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
			fall, ok := RainFallColor[colorCode]
			if !ok {
				fmt.Fprintf(os.Stderr, "Rainfall color %s is unknown\n", colorCode)
			}
			if rainFall < fall {
				rainFall = fall
				rainColor = c
			}
		}
	}
	return RainData{Time: n.Time, RainFall: rainFall, RainColor: rainColor}, nil
}

func (n NowcastImages) SeekRainInfo(pos PositionConfig) (RainInfo, error) {
	var res RainInfo
	for i := range n.Nowcast {
		rain, err := n.Nowcast[i].SeekRainData(pos)
		if err != nil {
			return res, fmt.Errorf("seek rain data : %w", err)
		}
		rain.ImageID = i
		res.Rains = append(res.Rains, rain)
	}
	return res, nil
}

func (r RainInfo) Rainy() bool {
	for i := range r.Rains {
		if r.Rains[i].RainFall > 0 {
			return true
		}
	}
	return false
}
