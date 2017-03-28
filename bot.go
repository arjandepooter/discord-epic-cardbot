package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/arjandepooter/discord-epic-cardbot/epicapi"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/bwmarrin/discordgo"
	"github.com/namsral/flag"
	"github.com/robfig/cron"
)

var (
	discord  *discordgo.Session
	schedule *cron.Cron
	cards    map[string]*epicapi.Card
	pattern  *regexp.Regexp
	index    bleve.Index
)

func getIndex(path string) (bleve.Index, error) {
	cardIndex, err := bleve.Open(path)

	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping := getIndexMapping()
		return bleve.New(path, indexMapping)
	} else if err != nil {
		return nil, err
	}

	return cardIndex, nil
}

func getIndexMapping() mapping.IndexMapping {
	indexMapping := bleve.NewIndexMapping()
	cardMapping := bleve.NewDocumentMapping()
	nameFieldMapping := bleve.NewTextFieldMapping()
	nameFieldMapping.Analyzer = "en"
	cardMapping.AddFieldMappingsAt("name", nameFieldMapping)
	indexMapping.DefaultMapping = cardMapping
	indexMapping.AddDocumentMapping("card", cardMapping)

	return indexMapping
}

func searchCard(cardName string) (*epicapi.Card, bool) {
	query := bleve.NewMatchQuery(cardName)
	query.SetField("name")
	search := bleve.NewSearchRequest(query)
	search.Size = 1
	result, err := index.Search(search)
	if err != nil {
		log.WithError(err).Error("Can't search index")
		return nil, false
	}
	if result.Total > 0 {
		card, found := cards[result.Hits[0].ID]
		return card, found
	}

	return nil, false
}

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
			card, found := searchCard(cardName)

			if found {
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
		card, found := searchCard(cardName)

		if found {
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
		err       error
		token     = flag.String("token", "", "Discord Authentication Token")
		indexName = flag.String("index", "epiccardindex", "Bleve index location")
	)

	flag.Parse()

	discord, err = discordgo.New(*token)
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

	index, err = getIndex(*indexName)
	if err != nil {
		log.WithError(err).Fatal("Error while opening index")
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
