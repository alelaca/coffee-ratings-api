package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alelaca/coffee-ratings-api/cmd/customerrors"

	"github.com/alelaca/coffee-ratings-api/cmd/services"
	"github.com/gin-gonic/gin"
)

type RatingDTO struct {
	CoffeeType string `json:"coffeeType"`
	StarRating string `json:"starRating"`
}

type Ratings struct {
	Service *services.Ratings
}

func InitializeRatingsController(service *services.Ratings) *Ratings {
	return &Ratings{
		Service: service,
	}
}

func (r *Ratings) CreateRating(c *gin.Context) {
	requestContent := c.Request.Body
	body, err := ioutil.ReadAll(requestContent)
	if err != nil {
		abortWithCustomError(c, http.StatusInternalServerError, err)
	}

	var ratingDTO RatingDTO
	if err := json.Unmarshal(body, &ratingDTO); err != nil {
		abortWithCustomError(c, http.StatusBadRequest, fmt.Errorf("invalid body in request"))
	}

	if ratingDTO.CoffeeType == "" {
		abortWithCustomError(c, http.StatusBadRequest, fmt.Errorf("coffeeType parameter is required"))
	}

	if ratingDTO.StarRating == "" {
		abortWithCustomError(c, http.StatusBadRequest, fmt.Errorf("starRating parameter is required"))
	}

	err = r.Service.CreateRating(ratingDTO.CoffeeType, ratingDTO.StarRating)
	if err != nil {
		abortWithCustomError(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, "Rating created")
}

func (r *Ratings) GetCoffeeTypeList(c *gin.Context) {
	coffeeTypeList, err := r.Service.GetCoffeeTypeList()
	if err != nil {
		abortWithCustomError(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, coffeeTypeList)
}

func (r *Ratings) GetRating(c *gin.Context) {
	coffeeType := c.Query("coffeeType")
	if coffeeType == "" {
		abortWithCustomError(c, http.StatusBadRequest, fmt.Errorf("coffeeType query parameter is required"))
	}

	rating, err := r.Service.GetRating(coffeeType)
	if err != nil {
		abortWithCustomError(c, http.StatusInternalServerError, err)
	}

	if rating == nil {
		c.JSON(http.StatusBadRequest, gin.H{"coffeeType": "not rated yet coffee type taken from query param"})
		return
	}

	ratingDTO := RatingDTO{
		CoffeeType: rating.CoffeeType,
		StarRating: fmt.Sprintf("%d", rating.Stars.Given) + "/" + fmt.Sprintf("%d", rating.Stars.Max),
	}

	c.JSON(http.StatusOK, ratingDTO)
}

func (r *Ratings) GetRecommendation(c *gin.Context) {
	rating, err := r.Service.GetRecommendation()
	if err != nil {
		abortWithCustomError(c, http.StatusInternalServerError, err)
		return
	}

	if rating == nil {
		c.JSON(http.StatusOK, gin.H{"message": "NO_RECOMMENDATIONS_AVAILABLE"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coffeeType": rating.CoffeeType})
}

func abortWithCustomError(c *gin.Context, defaultStatus int, err error) {
	if apiError, ok := err.(*customerrors.ApiError); ok && apiError.HttpStatusCode != 0 {
		c.AbortWithStatusJSON(apiError.HttpStatusCode, gin.H{"error": apiError.Error()})
	} else {
		c.AbortWithStatusJSON(defaultStatus, gin.H{"error": err.Error()})
	}
}
