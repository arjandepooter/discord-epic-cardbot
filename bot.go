package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/arjandepooter/discord-epic-cardbot/epicapi"
	"github.com/bwmarrin/discordgo"
	"github.com/namsral/flag"
	"github.com/robfig/cron"
)

var (
	discord  *discordgo.Session
	schedule *cron.Cron
	cards    map[string]*epicapi.Card
	pattern  *regexp.Regexp
)

func sendCard(card *epicapi.Card, session *discordgo.Session, channelID string) error {
	embed := new(discordgo.MessageEmbed)
	embed.Image = new(discordgo.MessageEmbedImage)
	embed.Image.URL = epicapi.BaseURL + card.ImageSource

	_, err := discord.ChannelMessageSendEmbed(channelID, embed)

	return err
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := m.Message.Content

	if pattern.MatchString(msg) {
		matches := pattern.FindAllStringSubmatch(msg, -1)

		for _, match := range matches {
			cardName := strings.ToLower(match[1])
			card, exists := cards[cardName]

			if exists {
				log.Info(fmt.Sprintf("Card found: %s", card.Name))
				err := sendCard(card, s, m.ChannelID)

				if err != nil {
					log.WithError(err).Error("Can't send card")
				}
			}
		}
	} else if strings.HasPrefix(msg, "!card ") {
		cardName := strings.ToLower(strings.TrimSpace(msg[6:]))
		log.Info(fmt.Sprintf("%s requested card `%s`", m.Author.Username, cardName))
		card, exists := cards[cardName]

		if exists {
			log.Info(fmt.Sprintf("Card found: %s", card.Name))
			err := sendCard(card, s, m.ChannelID)

			if err != nil {
				log.WithError(err).Error("Can't send card")
			}
		} else {
			log.Info(fmt.Sprintf("Card not found"))
			s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("Sorry, can't find a card with the name '%s' :cry:", cardName))
		}
	}
}

func init() {
	cards = make(map[string]*epicapi.Card)
	pattern = regexp.MustCompile(`\[{2}([^\]]+)\]{2}`)
}

func main() {
	var (
		err   error
		Token = flag.String("token", "", "Discord Authentication Token")
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
	schedule.AddFunc("@daily", updateCardDatabase)
	schedule.Start()

	log.Info("Epic CardBot is ready and waiting for commands")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
