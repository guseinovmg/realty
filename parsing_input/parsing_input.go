package parsing_input

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

/*
mistral.ai
create function that fill struct

	type GetAdvListRequest struct {
		Currency     string
		MinPrice     int64
		MaxPrice     int64
		MinLongitude float64
		MaxLongitude float64
		MinLatitude  float64
		MaxLatitude  float64
		CountryCode  string
		Location     string
		Page         int
		FirstNew     bool
	}

with values from request.URL.Query().Get.
If request.URL.Query().Get will return empty string, do nothing.
If parsing to int64, float64, bool returns error, then return the error with name of the field.
*/
func ParseQueryToGetAdvListRequest(query url.Values, req *dto.GetAdvListRequest) error {
	req.Currency = query.Get("currency")

	var value string

	value = query.Get("minPrice")
	if value != "" {
		minPrice, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("minPrice: %w", err)
		}
		req.MinPrice = minPrice
	}

	value = query.Get("maxPrice")
	if value != "" {
		maxPrice, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("maxPrice: %w", err)
		}
		req.MaxPrice = maxPrice
	}

	value = query.Get("minLongitude")
	if value != "" {
		minLongitude, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("minLongitude: %w", err)
		}
		req.MinLongitude = minLongitude
	}

	value = query.Get("maxLongitude")
	if value != "" {
		maxLongitude, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("maxLongitude: %w", err)
		}
		req.MaxLongitude = maxLongitude
	}

	value = query.Get("minLatitude")
	if value != "" {
		minLatitude, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("minLatitude: %w", err)
		}
		req.MinLatitude = minLatitude
	}

	value = query.Get("maxLatitude")
	if value != "" {
		maxLatitude, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("maxLatitude: %w", err)
		}
		req.MaxLatitude = maxLatitude
	}

	req.CountryCode = query.Get("countryCode")
	req.Location = query.Get("location")

	value = query.Get("page")
	if value != "" {
		page, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("page: %w", err)
		}
		req.Page = page
	}

	value = query.Get("firstNew")
	if value != "" {
		firstNew, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("firstNew: %w", err)
		}
		req.FirstNew = firstNew
	}

	return nil
}
