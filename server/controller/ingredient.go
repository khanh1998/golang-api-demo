package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uss-kelvin/golang-api/server/model"
	"go.mongodb.org/mongo-driver/bson"
)

type IngredientController struct {
	model model.Ingredients
}

func New(model model.Ingredients) IngredientController {
	return IngredientController{
		model: model,
	}
}

func (i *IngredientController) GetAll(c *gin.Context) {
	ingredients, err := i.model.GetAll()
	if err != nil {
		log.Panic(err)
	}
	c.IndentedJSON(http.StatusOK, ingredients)
}

func (i *IngredientController) GetById(c *gin.Context) {
	id := c.Param("id")
	ingredient, err := i.model.GetById(id)
	if err != nil {
		log.Panic(err)
	}
	c.IndentedJSON(http.StatusOK, ingredient)
}

func (i *IngredientController) AddOne(c *gin.Context) {
	var input model.Ingredient
	if err := c.BindJSON(&input); err != nil {
		log.Panic(err)
	}
	id, err := i.model.AddOne(input)
	if err != nil {
		log.Panic(err)
	}
	inserted, err := i.model.GetById(id)
	if err != nil {
		log.Panic(err)
	}
	c.IndentedJSON(http.StatusOK, inserted)
}

func (i *IngredientController) UpdateOne(c *gin.Context) {
	var input model.Ingredient
	id := c.Param("id")
	if err := c.BindJSON(&input); err != nil {
		log.Panic(err)
	}
	success, err := i.model.UpdateOne(id, input)
	if err != nil {
		log.Panic(err)
	}
	if success {
		ingredient, err := i.model.GetById(id)
		if err != nil {
			log.Panic(err)
		}
		c.IndentedJSON(http.StatusOK, ingredient)
	}
}

func (i *IngredientController) DeleteOne(c *gin.Context) {
	id := c.Param("id")
	success := i.model.DeleteOne(id)
	if success {
		c.IndentedJSON(http.StatusOK, bson.M{"message": fmt.Sprintf("Delete %v successfully", id)})
	} else {
		c.IndentedJSON(http.StatusOK, bson.M{"message": fmt.Sprintf("Delete %v fail", id)})
	}
}
