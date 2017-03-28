package main

import (
	"fmt"

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
		cards[card.Code] = card
		index.Index(card.Code, *card)
	}

	log.Info(fmt.Sprintf("%d cards found", len(cardList)))
}
