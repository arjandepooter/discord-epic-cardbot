package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/arjandepooter/discord-epic-cardbot/epicapi"
)

func updateCardDatabase() {
	log.Info("Updating card database")

	cards, err := epicapi.GetAllCards()
	if err != nil {
		log.WithError(err).Error("Error when fetching cards")
		return
	}

	log.Info(fmt.Sprintf("%d cards found", len(cards)))
}
