package main

import (
	"fmt"
	"log"

	"github.com/uss-kelvin/golang-api/server"
	"github.com/uss-kelvin/golang-api/server/config"
)

func main() {
	connection, err := config.GetConnection()
	if err != nil {
		log.Panic(err)
	}
	err = connection.TestConnection()
	if err != nil {
		log.Panic(err)
	}
	env, err := config.LoadEnv("")
	if err != nil {
		log.Panic(err)
	}
	app, err := server.NewServer(connection, env.DatabaseName)
	if err != nil {
		log.Panic(err)
	}
	if err = app.Start(env.Host); err != nil {
		log.Panic(err)
	}
	fmt.Printf("Server is running at %v \n", env.Host)
}
