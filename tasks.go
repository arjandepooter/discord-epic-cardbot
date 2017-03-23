package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/arjandepooter/discord-epic-cardbot/epicapi"
)

func updateCardDatabase() {
	log.Info("Updating card database")

	cardList, err := epicapi.GetAllCards()
	if err != nil {
		log.WithError(err).Error("Error when fetching cards")
		return
	}

	for _, card := range cardList {
		cards[strings.ToLower(card.Name)] = card
	}

	log.Info(fmt.Sprintf("%d cards found", len(cardList)))
}
