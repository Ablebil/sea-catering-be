package main

import (
	_ "github.com/Ablebil/sea-catering-be/docs"
	"github.com/Ablebil/sea-catering-be/internal/bootstrap"
)

// @title           Sea Catering API
// @version         1.0
// @description     This is the API documentation for the Sea Catering application.

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Enter your bearer token in the format `Bearer {token}`

func main() {
	if err := bootstrap.Start(); err != nil {
		panic(err)
	}
}
