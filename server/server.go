package server

import (
	"github.com/gin-gonic/gin"
	"github.com/uss-kelvin/golang-api/server/config"
	"github.com/uss-kelvin/golang-api/server/controller"
	"github.com/uss-kelvin/golang-api/server/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	router     *gin.Engine
	connection *config.Connection
	database   *mongo.Database
}

func NewServer(con *config.Connection, databaseName string) (*Server, error) {
	database := con.GetDatabase(databaseName)
	server := Server{
		connection: con,
		database:   database,
	}
	server.setupRouter()
	return &server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	ingredientModel := model.New("ingredients", s.database)
	ingredient := controller.New(*ingredientModel)
	router.GET("/ingredients", ingredient.GetAll)
	router.GET("/ingredients/:id", ingredient.GetById)
	router.POST("/ingredients", ingredient.AddOne)
	router.PUT("/ingredients/:id", ingredient.UpdateOne)
	router.DELETE("/ingredients/:id", ingredient.DeleteOne)
	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
