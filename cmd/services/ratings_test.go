package services

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alelaca/coffee-ratings-api/cmd/customerrors"
	"github.com/alelaca/coffee-ratings-api/cmd/dao"
	"github.com/alelaca/coffee-ratings-api/cmd/entities"
)

func TestStarsInvalidFormat(t *testing.T) {
	starsParam := "strinvalue"
	_, _, err := parseStars(starsParam)

	if err == nil {
		t.Fatalf("Test should have failed")
	}
}

func TestStarsNoMax(t *testing.T) {
	starsParam := "5/"
	_, _, err := parseStars(starsParam)

	if err == nil {
		t.Fatalf("Test should have failed")
	}
}

func TestStarsNoGiven(t *testing.T) {
	starsParam := "/5"
	_, _, err := parseStars(starsParam)

	if err == nil {
		t.Fatalf("Test should have failed")
	}
}

func TestStarsMaxOutOfRange(t *testing.T) {
	starsParam := "3/10"
	_, _, err := parseStars(starsParam)

	if err == nil {
		t.Fatalf("Test should have failed")
	}
}

func TestStarsGivenOutofRange(t *testing.T) {
	starsParam := "7/5"
	_, _, err := parseStars(starsParam)

	if err == nil {
		t.Fatalf("Test should have failed")
	}
}

func TestCreateRating_DAOError(t *testing.T) {
	daoMock := dao.RatingsRepositoryMock{
		NextError: fmt.Errorf("error!"),
	}

	ratingsService := Ratings{
		dao: daoMock,
	}

	err := ratingsService.CreateRating("coffee-type", "4/5")

	if err == nil {
		t.Fatalf("Test should have failed")
	}

	if apiError, ok := err.(customerrors.ApiError); ok {
		expectedErrorCode := http.StatusInternalServerError
		if apiError.HttpStatusCode != expectedErrorCode {
			t.Fatalf("Test failed. Expected error code %d, got %d", expectedErrorCode, apiError.HttpStatusCode)
		}
	} else {
		t.Fatalf("Expected an api error with a status code. Got %s", err.Error())
	}
}

func TestGetRating_DAOError(t *testing.T) {
	daoMock := dao.RatingsRepositoryMock{
		NextError: fmt.Errorf("error!"),
	}

	ratingsService := Ratings{
		dao: daoMock,
	}

	_, err := ratingsService.GetRating("coffee-type")

	if err == nil {
		t.Fatalf("Test should have failed")
	}

	if apiError, ok := err.(customerrors.ApiError); ok {
		expectedErrorCode := http.StatusInternalServerError
		if apiError.HttpStatusCode != expectedErrorCode {
			t.Fatalf("Test failed. Expected error code %d, got %d", expectedErrorCode, apiError.HttpStatusCode)
		}
	} else {
		t.Fatalf("Expected an api error with a status code. Got %s", err.Error())
	}
}

func TestGetRating_NotFound(t *testing.T) {
	daoMock := dao.RatingsRepositoryMock{
		NextRating: nil,
	}

	ratingsService := Ratings{
		dao: daoMock,
	}

	rating, err := ratingsService.GetRating("coffee-type")

	if err != nil {
		t.Fatalf("Test shouldnt have failed")
	}

	if rating != nil {
		t.Fatalf("rating should be nil")
	}
}

func TestGetRatingList_DAOError(t *testing.T) {
	daoMock := dao.RatingsRepositoryMock{
		NextError: fmt.Errorf("error!"),
	}

	ratingsService := Ratings{
		dao: daoMock,
	}

	_, err := ratingsService.GetCoffeeTypeList()

	if err == nil {
		t.Fatalf("Test should have failed")
	}

	if apiError, ok := err.(customerrors.ApiError); ok {
		expectedErrorCode := http.StatusInternalServerError
		if apiError.HttpStatusCode != expectedErrorCode {
			t.Fatalf("Test failed. Expected error code %d, got %d", expectedErrorCode, apiError.HttpStatusCode)
		}
	} else {
		t.Fatalf("Expected an api error with a status code. Got %s", err.Error())
	}
}

func TestGetRecommendation_NoRatings(t *testing.T) {
	ratingListResult := []entities.Rating{}

	daoMock := dao.RatingsRepositoryMock{
		NextRatingList: ratingListResult,
	}

	ratingsService := Ratings{
		dao: daoMock,
	}

	rating, err := ratingsService.GetRecommendation()

	if err != nil {
		t.Fatalf("Test shouldnt have failed")
	}

	if rating != nil {
		t.Fatalf("rating should be nil")
	}
}

func TestGetRecommendation_OldestRating(t *testing.T) {

	date := time.Now()

	ratingListResult := []entities.Rating{
		{
			CoffeeType: "2",
			Stars: entities.Stars{
				Given: 4,
				Max:   5,
			},
			DateCreated: date.Add(time.Hour),
		},
		{
			CoffeeType: "3",
			Stars: entities.Stars{
				Given: 4,
				Max:   5,
			},
			DateCreated: date.Add(2 * time.Hour),
		},
		{
			CoffeeType: "1",
			Stars: entities.Stars{
				Given: 4,
				Max:   5,
			},
			DateCreated: date.Add(3 * time.Hour),
		},
	}

	daoMock := dao.RatingsRepositoryMock{
		NextRatingList: ratingListResult,
	}

	ratingsService := Ratings{
		dao: daoMock,
	}

	rating, err := ratingsService.GetRecommendation()

	if err != nil {
		t.Fatalf("Test shouldnt have failed")
	}

	if rating == nil {
		t.Fatalf("rating shouldnt be nil")
	}

	expectedRatingCoffeeType := "2"
	if rating.CoffeeType != expectedRatingCoffeeType {
		t.Fatalf("Test failed. Coffee type expected %s, got %s", expectedRatingCoffeeType, rating.CoffeeType)
	}
}
