package dao

import (
	"github.com/alelaca/coffee-ratings-api/cmd/entities"
)

type Ratings interface {
	CreateRating(rating entities.Rating) error
	GetRatingList() ([]entities.Rating, error)
	GetRating(coffeeType string) (*entities.Rating, error)
}
