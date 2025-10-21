package main

import (
	"log"
	"rest"
	"rest/pkg/handler"
	"rest/pkg/repository"
	"rest/pkg/service"
)

func main() {

	repos := repository.NewRepository(nil) /* add db*/
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(rest.Server)

	if err := srv.Run("8081", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}
