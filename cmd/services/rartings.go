package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/alelaca/coffee-ratings-api/cmd/dao"

	"github.com/alelaca/coffee-ratings-api/cmd/customerrors"
	"github.com/alelaca/coffee-ratings-api/cmd/entities"
)

const (
	MaxStarsRating    = 5
	StarsRatingFormat = "given/max"
)

type Ratings struct {
	dao dao.Ratings
}

func InitializeRatingsService(dao dao.Ratings) *Ratings {
	return &Ratings{
		dao: dao,
	}
}

func (r *Ratings) CreateRating(coffeeType string, stars string) error {
	given, max, err := parseStars(stars)
	if err != nil {
		return err
	}

	newRating := entities.Rating{
		CoffeeType: coffeeType,
		Stars: entities.Stars{
			Given: given,
			Max:   max,
		},
		DateCreated: time.Now(),
	}

	err = r.dao.CreateRating(newRating)
	if err != nil {
		return customerrors.CreateApiError(http.StatusInternalServerError, fmt.Sprintf("error saving coffee rating, error: %s", err.Error()))
	}

	return nil
}

func (r *Ratings) GetCoffeeTypeList() ([]string, error) {
	ratingsList, err := r.dao.GetRatingList()
	if err != nil {
		return nil, customerrors.CreateApiError(http.StatusInternalServerError, fmt.Sprintf("error getting coffee rating list, error: %s", err.Error()))
	}

	coffeeTypeList := []string{}
	for _, v := range ratingsList {
		coffeeTypeList = append(coffeeTypeList, v.CoffeeType)
	}

	return coffeeTypeList, nil
}

func (r *Ratings) GetRating(coffeeType string) (*entities.Rating, error) {
	rating, err := r.dao.GetRating(coffeeType)
	if err != nil {
		return nil, customerrors.CreateApiError(http.StatusInternalServerError, fmt.Sprintf("error getting coffee rating '%s', error: %s", coffeeType, err.Error()))
	}

	return rating, nil
}

func (r *Ratings) GetRecommendation() (*entities.Rating, error) {
	ratings, err := r.dao.GetRatingList()
	if err != nil {
		return nil, customerrors.CreateApiError(http.StatusInternalServerError, fmt.Sprintf("error generating coffee recommendation, error: %s", err.Error()))
	}

	if len(ratings) == 0 {
		return nil, nil
	}

	oldestDate := ratings[0].DateCreated
	var oldest entities.Rating
	for _, v := range ratings {
		if v.Stars.Given > 3 {
			if v.DateCreated == oldestDate || v.DateCreated.Before(oldestDate) {
				oldestDate = v.DateCreated
				oldest = v
			}
		}
	}

	return &oldest, nil
}

func parseStars(stars string) (given, max int, err error) {
	given, max = 0, 0

	re := regexp.MustCompile("\\d+")
	regexResult := re.FindAllString(stars, -1)

	if len(regexResult) != 2 {
		err = customerrors.CreateApiError(http.StatusBadRequest, fmt.Sprintf("invalid stars rating format, received: '%s', expected: '%s'", stars, StarsRatingFormat))
		return
	}

	if regexResult[0]+"/"+regexResult[1] != stars {
		err = customerrors.CreateApiError(http.StatusBadRequest, fmt.Sprintf("invalid stars rating format, received: '%s', expected: '%s'", stars, StarsRatingFormat))
		return
	}

	given, error := strconv.Atoi(regexResult[0])
	if error != nil {
		err = customerrors.CreateApiError(http.StatusBadRequest, fmt.Sprintf("invalid stars rating format, given '%s' is not a number", regexResult[0]))
		return
	}
	max, error = strconv.Atoi(regexResult[1])
	if error != nil {
		err = customerrors.CreateApiError(http.StatusBadRequest, fmt.Sprintf("invalid stars rating format, max '%s' is not a number", regexResult[1]))
		return
	}

	if max > MaxStarsRating {
		err = customerrors.CreateApiError(http.StatusBadRequest, fmt.Sprintf("invalid max stars rating, received: '%d', max allowed: '%d'", max, MaxStarsRating))
		return
	}

	if given > max || given < 1 {
		err = customerrors.CreateApiError(http.StatusBadRequest, fmt.Sprintf("invalid given stars rating, given rating cannot be grater than max and less than 1"))
		return
	}

	return
}
