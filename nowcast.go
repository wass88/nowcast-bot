package main

import "fmt"

type PositionConfig struct {
	X            int
	Y            int
	CursorHeight int
	MapID        int
	Trim         TrimConfig
}
type TrimConfig struct {
	Trim   bool
	X      int
	Y      int
	Width  int
	Height int
}

type NowcastComposer struct {
	Pos       PositionConfig
	Images    NowcastImages
	Annotated NowcastImages
	Rain      RainInfo
}

func NewComposer(pos PositionConfig) (NowcastComposer, error) {
	var composer NowcastComposer
	images, err := DownloadableImage(pos.MapID)
	if err != nil {
		return composer, fmt.Errorf("downloadable image : %w", err)
	}
	err = images.Fetch()
	if err != nil {
		return composer, fmt.Errorf("fetch : %w", err)
	}
	rain, err := images.SeekRainInfo(pos)
	if err != nil {
		return composer, fmt.Errorf("seek rain info : %w", err)
	}
	charts, err := GenerateCharts(rain, pos.Trim.GetWidth())
	if err != nil {
		return composer, fmt.Errorf("generate charts : %w", err)
	}
	annotated, err := images.Annotate(pos, charts)
	if err != nil {
		return composer, fmt.Errorf("annotate : %w", err)
	}

	return NowcastComposer{Pos: pos, Images: images, Annotated: annotated, Rain: rain}, nil
}

func (c *NowcastComposer) Rainy() bool {
	return c.Rain.Rainy()
}
