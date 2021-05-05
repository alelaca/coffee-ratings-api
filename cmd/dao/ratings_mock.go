package dao

import "github.com/alelaca/coffee-ratings-api/cmd/entities"

type RatingsRepositoryMock struct {
	NextRatingList []entities.Rating
	NextRating     *entities.Rating
	NextError      error
}

func (r RatingsRepositoryMock) CreateRating(rating entities.Rating) error {
	return r.NextError
}

func (r RatingsRepositoryMock) GetRatingList() ([]entities.Rating, error) {
	return r.NextRatingList, r.NextError
}

func (r RatingsRepositoryMock) GetRating(coffeeType string) (*entities.Rating, error) {
	return r.NextRating, r.NextError
}
