package main

import (
	"Geonavigator-bot/coordinats"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var description = "Вас приветствует МСК Навигатор! Введите одну из координат участка (x, y) и вы получите точку на карте, к которой сможете доехать"

func main() {
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Text {
		case "/start":
			startMessage := tgbotapi.NewMessage(update.Message.Chat.ID, description)
			bot.Send(startMessage)
		default:
			mskCoord, err := coordinats.ParseMSKCoordinate(update.Message.Text)
			if err != nil {
				errorMess := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				bot.Send(errorMess)
				continue
			}
			wgsCoord, err := mskCoord.MSKToWGS()
			if err != nil {
				errorMess := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				bot.Send(errorMess)
				continue
			}
			location := tgbotapi.NewLocation(update.Message.Chat.ID, wgsCoord.Lat, wgsCoord.Long)
			bot.Send(location)
		}
	}
}
