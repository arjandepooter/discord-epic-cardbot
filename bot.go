package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/arjandepooter/discord-epic-cardbot/epicapi"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var (
	discord  *discordgo.Session
	schedule *cron.Cron
	cards    map[string]*epicapi.Card
)

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := m.Message.Content
	if strings.HasPrefix(msg, "!card ") {
		cardName := strings.ToLower(strings.TrimSpace(msg[6:]))
		log.Info(fmt.Sprintf("%s requested card `%s`", m.Author.Username, cardName))
		card, exists := cards[cardName]

		if exists {
			log.Info(fmt.Sprintf("Card found: %s", card.Name))
			embed := new(discordgo.MessageEmbed)
			embed.Image = new(discordgo.MessageEmbedImage)
			embed.Image.URL = epicapi.BaseURL + card.ImageSource

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		} else {
			log.Info(fmt.Sprintf("Card not found"))
			s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("Sorry, can't find a card with the name '%s' :cry:", cardName))
		}
	}
}

func main() {
	var (
		err   error
		Token = flag.String("t", "", "Discord Authentication Token")
	)

	flag.Parse()
	if len(*Token) == 0 {
		*Token = os.Getenv("TOKEN")
	}

	cards = make(map[string]*epicapi.Card)

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
