package dao

import (
	"github.com/alelaca/coffee-ratings-api/cmd/entities"
)

type LocalRepository struct {
	ratings map[string][]entities.Rating
}

func InitializeLocalRepository() *LocalRepository {
	return &LocalRepository{
		ratings: make(map[string][]entities.Rating),
	}
}

func (l *LocalRepository) CreateRating(rating entities.Rating) error {
	if _, ok := l.ratings[rating.CoffeeType]; !ok {
		l.ratings[rating.CoffeeType] = []entities.Rating{rating}
		return nil
	}

	l.ratings[rating.CoffeeType] = append(l.ratings[rating.CoffeeType], rating)
	return nil
}

func (l *LocalRepository) GetRatingList() ([]entities.Rating, error) {
	coffeeTypes := []entities.Rating{}
	for k := range l.ratings {
		coffeeTypes = append(coffeeTypes, l.ratings[k][len(l.ratings[k])-1]) // taking the list of most recent ratings
	}

	return coffeeTypes, nil
}

func (l *LocalRepository) GetRating(coffeeType string) (*entities.Rating, error) {
	if _, ok := l.ratings[coffeeType]; !ok {
		return nil, nil
	}

	mostRecentRatingIndex := len(l.ratings[coffeeType]) - 1
	return &l.ratings[coffeeType][mostRecentRatingIndex], nil
}
