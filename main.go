package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/gif"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	MapID         int    `envconfig:"MAP_ID" required:"true"`
	PosX          int    `envconfig:"POS_X" required:"true"`
	PosY          int    `envconfig:"POS_Y" required:"true"`
	SlackWebhook  string `envconfig:"SLACK_WEBHOOK" required:"true"`
	GyazoToken    string `envconfig:"GYAZO_TOKEN" required:"true"`
	Trim          string `envconfig:"TRIM"`
	RainThreshold int    `envconfig:"RAIN_THRESHOLD" default:"0"`
}

func main() {
	var env Env
	envconfig.MustProcess("", &env)
	pos := PositionConfig{
		X:            env.PosX,
		Y:            env.PosY,
		CursorHeight: 20,
		MapID:        env.MapID,
	}
	if env.Trim != "" {
		var trimConfig []int
		err := json.Unmarshal([]byte(env.Trim), &trimConfig)
		if err != nil || len(trimConfig) != 4 {
			panic(fmt.Errorf("env TRIM %s must be [x, y, w, h] : %w", env.Trim, err))
		}
		pos.Trim = TrimConfig{X: trimConfig[0], Y: trimConfig[1], Width: trimConfig[2], Height: trimConfig[3]}
		pos.Trim.Trim = true
	}
	fmt.Printf("Options %+v\n", pos)
	api := APIConfig{
		SlackWebhook: env.SlackWebhook,
		GyazoToken:   env.GyazoToken,
	}
	fmt.Printf("Fetch\n")
	composer, err := NewComposer(pos)
	if err != nil {
		panic(err)
	}
	if !composer.Rainy(env.RainThreshold) {
		fmt.Println("not rainy")
		return
	}
	fmt.Printf("Make GIF\n")
	img := composer.ComposeGif()
	file := &bytes.Buffer{}
	err = gif.EncodeAll(file, &img)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload GIF\n")
	url, err := UploadImage(file, api)
	if err != nil {
		panic(err)
	}

	msg := Message{
		Text:     fmt.Sprintf("%s", url),
		UserName: "雨情報",
		Emoji:    ":umbrella:",
	}

	fmt.Printf("Send to Slack\n")
	err = SendMessage(msg, api)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Completed\n")
}
