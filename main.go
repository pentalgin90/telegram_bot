package main

import (
	tgClient "bot/clients/telegram"
	"bot/consumer/event_consumer"
	"bot/events/telegram"
	"bot/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("server started")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

	//processor = processor.New(tgClient)

	//consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
