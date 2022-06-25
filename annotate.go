package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var nowcastWidth int = 940
var nowcastHeight int = 783

func DrawFont(img *image.RGBA, text string, x, y float64) {
	point := fixed.Point26_6{X: fixed.I(int(x)), Y: fixed.I(int(y))}
	col := color.Black

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

func GenerateCharts(r RainInfo, widthInt int) ([]image.Image, error) {
	width := (float64)(widthInt)
	height := 50.
	boxTop := 14.
	boxLeft := 10.
	boxesWidth := width - boxLeft*2
	boxHeight := 22.
	lineTop := 5.
	textTop := 11.

	chart := image.NewRGBA(image.Rect(0, 0, nowcastWidth, nowcastHeight))
	gc := draw2dimg.NewGraphicContext(chart)
	gc.SetFillColor(color.Gray{0xa8})
	draw2dkit.Rectangle(gc, 0, 0, width, height)
	gc.Fill()

	firstTime := r.Rains[0].Time
	lastTime := r.Rains[len(r.Rains)-1].Time
	boxTick := boxesWidth / lastTime.Sub(firstTime).Minutes()
	timeX := func(t time.Time) float64 {
		return boxTick*t.Sub(firstTime).Minutes() + boxLeft
	}

	// Draw Rain Box
	for i := range r.Rains {
		if i == len(r.Rains)-1 {
			continue
		}
		rain := r.Rains[i]
		if rain.RainFall == 0 {
			continue
		}
		nextTime := r.Rains[i+1].Time
		gc.SetFillColor(rain.RainColor)
		draw2dkit.Rectangle(gc, timeX(rain.Time), boxTop, timeX(nextTime), boxTop+boxHeight)
		gc.Fill()
	}

	// Draw Hour Line and Label
	nextHour := firstTime.Add(-time.Second).Truncate(time.Hour).Add(time.Hour)
	for h := nextHour; h.Before(lastTime); h = h.Add(time.Hour) {
		gc.SetFillColor(color.Gray{0x0a})
		lineX := timeX(h)
		thick := 1.
		if h.Add(9*time.Hour).Hour()%6 == 0 {
			thick = 2.
		}
		draw2dkit.Rectangle(gc, lineX, boxTop+lineTop, lineX+thick, boxTop+boxHeight-lineTop)
		gc.Fill()

		if h == nextHour {
			str := h.Add(9 * time.Hour).Format("01-02")
			DrawFont(chart, str, boxLeft, textTop)
		}
		str := h.Add(9 * time.Hour).Format("15")
		DrawFont(chart, str, lineX-5, boxTop+boxHeight+10)
	}
	gc.Close()

	// Add Cursor
	var charts []image.Image
	for i := range r.Rains {
		c := image.NewRGBA(image.Rect(0, 0, nowcastWidth, nowcastHeight))
		draw.Copy(c, image.Point{}, chart, chart.Bounds(), draw.Over, &draw.Options{})
		gc := draw2dimg.NewGraphicContext(c)
		lineX := timeX(r.Rains[i].Time)
		gc.SetFillColor(color.RGBA{0xff, 0x66, 0x22, 0xff})
		draw2dkit.Rectangle(gc, lineX-2, boxTop-2, lineX+2, boxTop+boxHeight+2)
		gc.Fill()
		charts = append(charts, c)
		gc.Close()
	}

	return charts, nil
}

func (img NowcastImage) AddCursor(pos PositionConfig, trim TrimConfig) (NowcastImage, error) {
	res := img.Image.(*image.RGBA)
	gc := draw2dimg.NewGraphicContext(res)
	gc.SetFillColor(color.RGBA{0xff, 0x66, 0x22, 0xff})
	x := float64(pos.X - trim.GetX())
	y := float64(pos.Y - trim.GetY())
	h := float64(pos.CursorHeight)
	draw2dkit.Rectangle(gc, x, y+1, x+1, y+h)
	gc.Fill()
	draw2dkit.Rectangle(gc, x-h/2+1, y+h/2, x+h/2, y+h/2+1)
	gc.Fill()
	img.Image = res
	return img, nil
}

func (img NowcastImage) AddChart(chart image.Image, trim TrimConfig) (NowcastImage, error) {
	res := image.NewRGBA(trim.GetSize())
	draw.Copy(res, image.Point{}, img.Image, trim.GetSize(), draw.Over, &draw.Options{})
	draw.Copy(res, image.Point{}, chart, trim.GetSize(), draw.Over, &draw.Options{})
	img.Image = res
	return img, nil
}

func (img NowcastImage) AddOnMap(mapImage image.Image, trim TrimConfig) (NowcastImage, error) {
	dx := trim.GetX()
	dy := trim.GetY()
	width := trim.GetWidth()
	height := trim.GetHeight()
	res := image.NewRGBA(trim.GetSize())
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			res.Set(x, y, color.Gray{0x30})
		}
	}
	draw.Copy(res, image.Point{}, mapImage, trim.GetBound(), draw.Over, &draw.Options{})
	transparented := img.Image
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, a := transparented.At(x+dx, y+dy).RGBA()
			if a > 0 {
				R, G, B, _ := res.At(x, y).RGBA()
				alphaF := uint32(3)
				alphaT := uint32(1)
				alpha := alphaF + alphaT
				r = r/alpha*alphaF + R/alpha*alphaT
				g = g/alpha*alphaF + G/alpha*alphaT
				b = b/alpha*alphaF + B/alpha*alphaT
				res.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xff})
			}
		}
	}
	img.Image = res
	return img, nil
}

func (img NowcastImage) Annotate(pos PositionConfig, chart image.Image, mapImg image.Image) (NowcastImage, error) {
	img, err := img.AddOnMap(mapImg, pos.Trim)
	if err != nil {
		return img, fmt.Errorf("add chart : %w", err)
	}
	img, err = img.AddCursor(pos, pos.Trim)
	if err != nil {
		return img, fmt.Errorf("add cursor : %w", err)
	}
	img, err = img.AddChart(chart, pos.Trim)
	if err != nil {
		return img, fmt.Errorf("add chart : %w", err)
	}
	return img, err
}

func (imgs NowcastImages) Annotate(pos PositionConfig, charts []image.Image) (NowcastImages, error) {
	for i := range imgs.Nowcast {
		img, err := imgs.Nowcast[i].Annotate(pos, charts[i], imgs.Map.Image)
		if err != nil {
			return imgs, fmt.Errorf("annotate : %w", err)
		}
		imgs.Nowcast[i] = img
	}
	return imgs, nil
}

func (t *TrimConfig) SetDefault() {
	if !t.Trim {
		t.X = 0
		t.Y = 0
		t.Width = nowcastWidth
		t.Height = nowcastHeight
	}
}

func (t TrimConfig) GetX() int {
	t.SetDefault()
	return t.X
}

func (t TrimConfig) GetY() int {
	t.SetDefault()
	return t.Y
}

func (t TrimConfig) GetWidth() int {
	t.SetDefault()
	return t.Width
}

func (t TrimConfig) GetHeight() int {
	t.SetDefault()
	return t.Height
}

func (t TrimConfig) GetBound() image.Rectangle {
	t.SetDefault()
	return image.Rect(t.X, t.Y, t.X+t.Width, t.Y+t.Height)
}

func (t TrimConfig) GetSize() image.Rectangle {
	t.SetDefault()
	return image.Rect(0, 0, t.Width, t.Height)
}
