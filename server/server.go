package server

import (
	"github.com/gin-gonic/gin"
	"github.com/uss-kelvin/golang-api/server/auth"
	"github.com/uss-kelvin/golang-api/server/config"
	"github.com/uss-kelvin/golang-api/server/controller"
	"github.com/uss-kelvin/golang-api/server/middleware"
	"github.com/uss-kelvin/golang-api/server/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	router     *gin.Engine
	connection *config.Connection
	database   *mongo.Database
	tokenMaker auth.Maker
	env        *config.Env
}

func NewServer(con *config.Connection, env *config.Env) (*Server, error) {
	database := con.GetDatabase(env.DatabaseName)
	tokenMaker, err := auth.NewPasetoMaker(env.SymmetricKey)
	if err != nil {
		return nil, err
	}
	server := Server{
		connection: con,
		database:   database,
		tokenMaker: tokenMaker,
		env:        env,
	}
	server.setupRouter()
	return &server, nil
}

func (s *Server) setupRouter() {
	authMiddleware := middleware.NewAuthMiddleware(s.tokenMaker)
	router := gin.Default()
	authRouter := router.Group("/").Use(authMiddleware)

	ingredientModel := model.NewIngredients("ingredients", s.database)
	ingredient := controller.New(*ingredientModel)
	authRouter.GET("/ingredients", ingredient.GetAll)
	authRouter.GET("/ingredients/:id", ingredient.GetById)
	authRouter.POST("/ingredients", ingredient.AddOne)
	authRouter.PUT("/ingredients/:id", ingredient.UpdateOne)
	authRouter.DELETE("/ingredients/:id", ingredient.DeleteOne)

	userModel := model.NewUsers("users", s.database)
	userController := controller.NewUserController(*userModel)
	authRouter.GET("/users/:id", userController.FindOne)
	router.POST("/users", userController.InsertOne)

	router.POST("/auth", userController.CreateLogin(s.tokenMaker))
	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
