package validator

import (
	"errors"
	"fmt"
	"realty/dto"
	"regexp"
	"time"
)

func IsValidUnixMicroId(id int64) bool {
	return id > 1720060451151465 && id < time.Now().UnixMicro()
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateLoginRequest(req *dto.LoginRequest) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	if err := validatePassword(req.Password); err != nil {
		return err
	}
	return nil
}

func ValidateRegisterRequest(req *dto.RegisterRequest) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	if err := validateName(req.Name); err != nil {
		return err
	}
	if err := validatePassword(req.Password); err != nil {
		return err
	}
	if err := validateInviteId(req.InviteId); err != nil {
		return err
	}
	return nil
}

func ValidateCreateAdvRequest(req *dto.CreateAdvRequest) error {
	if err := validateOriginLang(req.OriginLang); err != nil {
		return err
	}
	if err := validateTranslatedBy(req.TranslatedBy); err != nil {
		return err
	}
	if err := validateTranslatedTo(req.TranslatedTo); err != nil {
		return err
	}
	if err := validateTitle(req.Title); err != nil {
		return err
	}
	if err := validateDescription(req.Description); err != nil {
		return err
	}
	if err := validatePhotos(req.Photos); err != nil {
		return err
	}
	if err := validatePrice(req.Price); err != nil {
		return err
	}
	if err := validateCurrency(req.Currency); err != nil {
		return err
	}
	if err := validateCountry(req.Country); err != nil {
		return err
	}
	if err := validateCity(req.City); err != nil {
		return err
	}
	if err := validateAddress(req.Address); err != nil {
		return err
	}
	if err := validateLatitude(req.Latitude); err != nil {
		return err
	}
	if err := validateLongitude(req.Longitude); err != nil {
		return err
	}
	if err := validateUserComment(req.UserComment); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateAdvRequest(req *dto.UpdateAdvRequest) error {
	if err := validateOriginLang(req.OriginLang); err != nil {
		return err
	}
	if err := validateTranslatedBy(req.TranslatedBy); err != nil {
		return err
	}
	if err := validateTranslatedTo(req.TranslatedTo); err != nil {
		return err
	}
	if err := validateTitle(req.Title); err != nil {
		return err
	}
	if err := validateDescription(req.Description); err != nil {
		return err
	}
	if err := validatePhotos(req.Photos); err != nil {
		return err
	}
	if err := validatePrice(req.Price); err != nil {
		return err
	}
	if err := validateCurrency(req.Currency); err != nil {
		return err
	}
	if err := validateCountry(req.Country); err != nil {
		return err
	}
	if err := validateCity(req.City); err != nil {
		return err
	}
	if err := validateAddress(req.Address); err != nil {
		return err
	}
	if err := validateLatitude(req.Latitude); err != nil {
		return err
	}
	if err := validateLongitude(req.Longitude); err != nil {
		return err
	}
	if err := validateUserComment(req.UserComment); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateUserRequest(req *dto.UpdateUserRequest) error {
	if err := validateName(req.Name); err != nil {
		return err
	}
	if err := validateDescription(req.Description); err != nil {
		return err
	}
	return nil
}

func ValidateUpdatePasswordRequest(req *dto.UpdatePasswordRequest) error {
	if err := validatePassword(req.OldPassword); err != nil {
		return err
	}
	if err := validatePassword(req.NewPassword); err != nil {
		return err
	}
	return nil
}

func ValidateGetAdvListRequest(req *dto.GetAdvListRequest) error {
	if err := validateCurrency(req.Currency); err != nil {
		return fmt.Errorf("currency: %w", err)
	}
	if err := validateMinPrice(req.MinPrice); err != nil {
		return fmt.Errorf("minPrice: %w", err)
	}
	if err := validateMaxPrice(req.MaxPrice); err != nil {
		return fmt.Errorf("maxPrice: %w", err)
	}
	if err := validateMinLongitude(req.MinLongitude); err != nil {
		return fmt.Errorf("minLongitude: %w", err)
	}
	if err := validateMaxLongitude(req.MaxLongitude); err != nil {
		return fmt.Errorf("maxLongitude: %w", err)
	}
	if err := validateMinLatitude(req.MinLatitude); err != nil {
		return fmt.Errorf("minLatitude: %w", err)
	}
	if err := validateMaxLatitude(req.MaxLatitude); err != nil {
		return fmt.Errorf("maxLatitude: %w", err)
	}
	if err := validateCountryCode(req.CountryCode); err != nil {
		return fmt.Errorf("countryCode: %w", err)
	}
	if err := validateLocation(req.Location); err != nil {
		return fmt.Errorf("location: %w", err)
	}
	if err := validatePage(req.Page); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	return nil
}

func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 100 {
		return errors.New("invalid password")
	}
	return nil
}

func validateName(name string) error {
	if len(name) > 100 {
		return errors.New("invalid name")
	}
	return nil
}

func validateInviteId(inviteId string) error {
	if len(inviteId) != 0 && len(inviteId) > 20 {
		return errors.New("invalid inviteId")
	}
	return nil
}

func validateOriginLang(originLang int8) error {
	if originLang < -128 || originLang > 127 {
		return errors.New("invalid originLang")
	}
	return nil
}

func validateTranslatedBy(translatedBy int8) error {
	if translatedBy < -128 || translatedBy > 127 {
		return errors.New("invalid translatedBy")
	}
	return nil
}

func validateTranslatedTo(translatedTo string) error {
	if len(translatedTo) > 255 {
		return errors.New("invalid translatedTo")
	}
	return nil
}

func validateTitle(title string) error {
	if len(title) > 255 {
		return errors.New("invalid title")
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > 1000 {
		return errors.New("invalid description")
	}
	return nil
}

func validatePhotos(photos string) error {
	if len(photos) > 1000 {
		return errors.New("invalid photos")
	}
	return nil
}

func validatePrice(price int64) error {
	if price < 0 {
		return errors.New("invalid price")
	}
	return nil
}

func validateCountry(country string) error {
	if len(country) > 100 {
		return errors.New("invalid country")
	}
	return nil
}

func validateCity(city string) error {
	if len(city) > 100 {
		return errors.New("invalid city")
	}
	return nil
}

func validateAddress(address string) error {
	if len(address) > 255 {
		return errors.New("invalid address")
	}
	return nil
}

func validateLatitude(latitude float64) error {
	if latitude < -90 || latitude > 90 {
		return errors.New("invalid latitude")
	}
	return nil
}

func validateLongitude(longitude float64) error {
	if longitude < -180 || longitude > 180 {
		return errors.New("invalid longitude")
	}
	return nil
}

func validateUserComment(userComment string) error {
	if len(userComment) > 255 {
		return errors.New("invalid userComment")
	}
	return nil
}

func validateCurrency(currency string) error {
	if len(currency) != 3 {
		return fmt.Errorf("currency must be 3 characters long")
	}
	return nil
}

func validateMinPrice(minPrice int64) error {
	if minPrice < 0 {
		return fmt.Errorf("minPrice must be greater than or equal to 0")
	}
	return nil
}

func validateMaxPrice(maxPrice int64) error {
	if maxPrice < 0 {
		return fmt.Errorf("maxPrice must be greater than or equal to 0")
	}
	return nil
}

func validateMinLongitude(minLongitude float64) error {
	if minLongitude < -180 {
		return fmt.Errorf("minLongitude must be greater than or equal to -180")
	}
	if minLongitude > 180 {
		return fmt.Errorf("minLongitude must be less than or equal to 180")
	}
	return nil
}

func validateMaxLongitude(maxLongitude float64) error {
	if maxLongitude < -180 {
		return fmt.Errorf("maxLongitude must be greater than or equal to -180")
	}
	if maxLongitude > 180 {
		return fmt.Errorf("maxLongitude must be less than or equal to 180")
	}
	return nil
}

func validateMinLatitude(minLatitude float64) error {
	if minLatitude < -90 {
		return fmt.Errorf("minLatitude must be greater than or equal to -90")
	}
	if minLatitude > 90 {
		return fmt.Errorf("minLatitude must be less than or equal to 90")
	}
	return nil
}

func validateMaxLatitude(maxLatitude float64) error {
	if maxLatitude < -90 {
		return fmt.Errorf("maxLatitude must be greater than or equal to -90")
	}
	if maxLatitude > 90 {
		return fmt.Errorf("maxLatitude must be less than or equal to 90")
	}
	return nil
}

func validateCountryCode(countryCode string) error {
	if countryCode == "" {
		return nil
	}
	if len(countryCode) != 2 {
		return fmt.Errorf("countryCode must be 2 characters long")
	}
	return nil
}

func validateLocation(location string) error {
	if location == "" {
		return nil
	}
	if len(location) < 2 || len(location) > 255 {
		return fmt.Errorf("location must be between 2 and 255 characters long")
	}
	return nil
}

func validatePage(page int) error {
	if page < 1 {
		return fmt.Errorf("page must be greater than or equal to 1")
	}
	return nil
}
