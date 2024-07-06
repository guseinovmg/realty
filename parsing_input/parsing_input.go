package parsing_input

import (
	"encoding/json"
	"errors"
	"net/http"
	"realty/dto"
	"strconv"
	"strings"
)

func ParseRawJson(r *http.Request, m any) error {
	if r.Body != nil {
		ct := r.Header.Get("Content-Type")
		if ct != "" && strings.HasPrefix(ct, "application/json") {
			return json.NewDecoder(r.Body).Decode(&m)
		} else {
			return errors.New("Content-Type must be application/json")
		}
	} else {
		return errors.New("body is empty")
	}
}

func parsePostFormToUpdateAdvRequest(r *http.Request) (*dto.UpdateAdvRequest, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	updateAdvRequest := &dto.UpdateAdvRequest{}

	originLang, err := strconv.ParseInt(r.PostForm.Get("originLang"), 10, 8)
	if err == nil {
		updateAdvRequest.OriginLang = int8(originLang)
	}

	translatedBy, err := strconv.ParseInt(r.PostForm.Get("translatedBy"), 10, 8)
	if err == nil {
		updateAdvRequest.TranslatedBy = int8(translatedBy)
	}

	updateAdvRequest.TranslatedTo = r.PostForm.Get("translatedTo")
	updateAdvRequest.Title = r.PostForm.Get("title")
	updateAdvRequest.Description = r.PostForm.Get("description")
	updateAdvRequest.Photos = r.PostForm.Get("photos")

	price, err := strconv.ParseInt(r.PostForm.Get("price"), 10, 64)
	if err == nil {
		updateAdvRequest.Price = price
	}

	updateAdvRequest.Currency = r.PostForm.Get("currency")
	updateAdvRequest.Country = r.PostForm.Get("country")
	updateAdvRequest.City = r.PostForm.Get("city")
	updateAdvRequest.Address = r.PostForm.Get("address")

	latitude, err := strconv.ParseFloat(r.PostForm.Get("latitude"), 64)
	if err == nil {
		updateAdvRequest.Latitude = latitude
	}

	longitude, err := strconv.ParseFloat(r.PostForm.Get("longitude"), 64)
	if err == nil {
		updateAdvRequest.Longitude = longitude
	}

	updateAdvRequest.UserComment = r.PostForm.Get("userComment")

	return updateAdvRequest, nil
}
