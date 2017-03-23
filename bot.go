package main

import (
	"flag"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var (
	discord  *discordgo.Session
	schedule *cron.Cron
)

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info(m.Author.Username)
}

func main() {
	var (
		err   error
		Token = flag.String("t", "", "Discord Authentication Token")
	)
	flag.Parse()

	discord, err = discordgo.New(*Token)
	if err != nil {
		log.WithError(err).Fatal("Can't connect to Discord")
		return
	}

	discord.AddHandler(onMessage)

	err = discord.Open()
	if err != nil {
		log.WithError(err).Fatal("Failed to create Discord websocket")
		return
	}

	updateCardDatabase()
	schedule = cron.New()
	schedule.AddFunc("@hourly", updateCardDatabase)
	schedule.Start()

	log.Info("Epic CardBot is ready and waiting for commands")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
