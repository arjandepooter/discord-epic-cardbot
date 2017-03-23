package epicapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

// BaseURL is the root of the Epic card database
const BaseURL = "http://decks.epiccardgame.com"

// Card represents a card in an API response
type Card struct {
	PackCode      string `json:"pack_code"`
	PackName      string `json:"pack_name"`
	TypeCode      string `json:"type_code"`
	TypeName      string `json:"type_name"`
	FactionCode   string `json:"faction_code"`
	FactionName   string `json:"faction_name"`
	Position      int    `json:"position"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Cost          int    `json:"cost"`
	Text          string `json:"text"`
	Quantity      int    `json:"quantity"`
	CubeMaxCopies int    `json:"cube_max_copies"`
	Illustrator   string `json:"illustrator"`
	URL           string `json:"url"`
	ImageSource   string `json:"imagesrc"`
}

func makeRequest(path string) (response *http.Response, err error) {
	response, err = http.Get(BaseURL + path)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	return
}

// GetAllCards returns all cards from the API
func GetAllCards() (cards []*Card, err error) {
	cards = []*Card{}
	response, err := makeRequest("/api/public/cards/")

	if err != nil {
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&cards)

	return
}
