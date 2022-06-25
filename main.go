package main

import (
	"bytes"
	"fmt"
	"image/gif"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	MapID        int    `envconfig:"MAP_ID" required:"true"`
	PosX         int    `envconfig:"POS_X" required:"true"`
	PosY         int    `envconfig:"POS_Y" required:"true"`
	SlackWebhook string `envconfig:"SLACK_WEBHOOK" required:"true"`
	GyazoToken   string `envconfig:"GYAZO_TOKEN" required:"true"`
}

func main() {
	var env Env
	envconfig.MustProcess("", &env)
	pos := PositionConfig{
		X:            env.PosX,
		Y:            env.PosY,
		CursorHeight: 10,
		MapID:        env.MapID,
	}
	api := APIConfig{
		SlackWebhook: env.SlackWebhook,
		GyazoToken:   env.GyazoToken,
	}
	composer, err := NewComposer(pos)
	if err != nil {
		panic(err)
	}
	if !composer.Rainy() {
		fmt.Println("not rainy")
		return
	}
	img := composer.ComposeGif()
	file := &bytes.Buffer{}
	err = gif.EncodeAll(file, &img)
	if err != nil {
		panic(err)
	}
	url, err := UploadImage(file, api)
	if err != nil {
		panic(err)
	}

	msg := Message{
		Text:     fmt.Sprintf("%s", url),
		UserName: "雨情報",
		Emoji:    ":umbrella:",
	}

	err = SendMessage(msg, api)
	if err != nil {
		panic(err)
	}
}
