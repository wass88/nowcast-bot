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
	RainRange float64
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
	fallRange := 0.
	for i := 0; i < 2*pos.CursorHeight; i++ {
		dy := i - pos.CursorHeight - pos.CursorHeight/2
		dx := 0
		if i < pos.CursorHeight {
			dy = 0
			dx = i - pos.CursorHeight/2
		}
		c := n.Image.At(pos.X+dx, pos.Y+dy)
		r, g, b, a := c.RGBA()
		if a > 0 {
			fallRange += 1
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
	return RainData{Time: n.Time, RainFall: rainFall, RainColor: rainColor, RainRange: fallRange / float64(2*pos.CursorHeight)}, nil
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

func (r RainInfo) RainTotal() int {
	total := 0
	for i := range r.Rains {
		total += r.Rains[i].RainFall
	}
	return total
}
