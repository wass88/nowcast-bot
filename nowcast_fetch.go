package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type NowcastImage struct {
	Image image.Image `json:"-"`
	Time  time.Time   `json:"date"`
	Url   string      `json:"url"`
}

type NowcastImages struct {
	Nowcast []NowcastImage `json:"nowcast"`
	Map     NowcastImage   `json:"map"`
}

type RainDataTimes struct {
	ra    time.Time
	srf   time.Time
	srf15 time.Time
}

func FetchTime(url string) (time.Time, error) {
	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}, fmt.Errorf("get : %w", err)
	}
	defer resp.Body.Close()
	var res struct {
		Time time.Time `json:"time"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, fmt.Errorf("read all : %w", err)
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return time.Time{}, fmt.Errorf("unmarshal : %w", err)
	}
	return res.Time, nil
}

func FetchTimes() (RainDataTimes, error) {
	ra, err := FetchTime("https://www.jma.go.jp/bosai/rain/data/ra/time.json")
	if err != nil {
		return RainDataTimes{}, fmt.Errorf("fetch ra : %w", err)
	}
	srf, err := FetchTime("https://www.jma.go.jp/bosai/rain/data/srf/time.json")
	if err != nil {
		return RainDataTimes{}, fmt.Errorf("fetch srf : %w", err)
	}
	srf15, err := FetchTime("https://www.jma.go.jp/bosai/rain/data/srf15/time.json")
	if err != nil {
		return RainDataTimes{}, fmt.Errorf("fetch srf15 : %w", err)
	}
	return RainDataTimes{ra, srf, srf15}, nil
}

func DownloadableImage(mapID int) (NowcastImages, error) {
	var res NowcastImages
	times, err := FetchTimes()
	if err != nil {
		return res, fmt.Errorf("fetch times : %w", err)
	}
	images := []NowcastImage{}

	forcastTime := times.ra.Format("20060102150405")
	images = append(images, NowcastImage{
		Image: nil,
		Time:  times.ra,
		Url:   fmt.Sprintf("https://www.jma.go.jp/bosai/rain/data/ra/%s/rain01_%s_f%02d_a%02d.png", forcastTime, forcastTime, 0, mapID),
	})
	forcastTime = times.srf.Format("20060102150405")
	for hour := 1; hour <= 6; hour++ {
		t := times.srf.Add(time.Hour * time.Duration(hour))
		images = append(images, NowcastImage{
			Image: nil,
			Time:  t,
			Url:   fmt.Sprintf("https://www.jma.go.jp/bosai/rain/data/srf/%s/rain01_%s_f%02d_a%02d.png", forcastTime, forcastTime, hour, mapID),
		})
	}
	forcastTime = times.srf15.Format("20060102150405")
	for hour := 7; hour <= 15; hour++ {
		t := times.srf.Add(time.Hour * time.Duration(hour))
		images = append(images, NowcastImage{
			Image: nil,
			Time:  t,
			Url:   fmt.Sprintf("https://www.jma.go.jp/bosai/rain/data/srf15/%s/rain01_%s_f%02d_a%02d.png", forcastTime, forcastTime, hour, mapID),
		})
	}

	res.Nowcast = images
	res.Map = NowcastImage{
		Image: nil,
		Time:  time.Time{},
		Url:   fmt.Sprintf("https://www.jma.go.jp/bosai/rain/const/map/map_a%02d.png", mapID),
	}
	return res, nil
}

func FetchImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get : %w", err)
	}
	defer resp.Body.Close()
	return png.Decode(resp.Body)
}

func (n *NowcastImages) Fetch() error {
	for i := range n.Nowcast {
		image, err := FetchImage(n.Nowcast[i].Url)
		if err != nil {
			return fmt.Errorf("fetch image : %w", err)
		}
		n.Nowcast[i].Image = image
	}

	image, err := FetchImage(n.Map.Url)
	if err != nil {
		return fmt.Errorf("fetch image : %w", err)
	}
	n.Map.Image = image
	return nil
}

func DumpImage(image image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file : %w", err)
	}
	defer file.Close()
	return png.Encode(file, image)
}

func LoadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file : %w", err)
	}
	defer file.Close()
	image, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode file : %w", err)
	}
	return image, nil
}

func (n *NowcastImages) Dump(dir string) error {
	bytes, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("marshal : %w", err)
	}

	ioutil.WriteFile(dir+"/nowcast.json", bytes, 0644)

	for i := range n.Nowcast {
		filename := fmt.Sprintf("%s/%02di.png", dir, i)
		err := DumpImage(n.Nowcast[i].Image, filename)
		if err != nil {
			return fmt.Errorf("dump image %s : %w", filename, err)
		}
	}
	filename := fmt.Sprintf("%s/map.png", dir)
	err = DumpImage(n.Map.Image, filename)
	if err != nil {
		return fmt.Errorf("dump image %s : %w", filename, err)
	}
	return nil
}

func LoadNowcastImages(dir string) (NowcastImages, error) {
	var res NowcastImages
	file, err := os.Open(dir + "/nowcast.json")
	if err != nil {
		return res, fmt.Errorf("open file : %w", err)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return res, fmt.Errorf("read file : %w", err)
	}
	json.Unmarshal(bytes, &res)

	for i := range res.Nowcast {
		filename := fmt.Sprintf("%s/%02di.png", dir, i)
		image, err := LoadImage(filename)
		if err != nil {
			return res, fmt.Errorf("load image %s : %w", filename, err)
		}
		res.Nowcast[i].Image = image
	}

	filename := fmt.Sprintf("%s/map.png", dir)
	image, err := LoadImage(filename)
	if err != nil {
		return res, fmt.Errorf("load image %s : %w", filename, err)
	}
	res.Map.Image = image

	return res, nil
}
