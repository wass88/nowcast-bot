package main

import "image"

type Message struct {
}

type SlackConfig struct {
	SlackWebhook string
	GyazoToken   string
}

func SendMessage(img image.Image, config SlackConfig) error {
	//TODO
	return nil
}

func UploadImage(img image.Image, config SlackConfig) (string, error) {
	//TODO
	return "", nil
}
