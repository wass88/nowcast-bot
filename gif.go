package main

import (
	"image"
	"image/gif"

	"github.com/soniakeys/quant/median"
)

func (c NowcastComposer) ComposeGif() gif.GIF {
	delay := 100
	res := gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
	}
	for i, img := range c.Annotated.Nowcast {
		q := median.Quantizer(256)
		paletted := q.Paletted(img.Image)
		res.Image = append(res.Image, paletted)
		if c.Rain.Rains[i].RainFall > 0 {
			res.Delay = append(res.Delay, delay)
		} else {
			res.Delay = append(res.Delay, delay/5)
		}
	}
	return res
}
