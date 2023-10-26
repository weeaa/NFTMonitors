package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weeaa/nft/api"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/internal/services"
	"os"
)

func main() {
	c := make(chan struct{})
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pg, err := db.New(context.Background(), fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", os.Getenv("PSQL_USERNAME"), os.Getenv("PSQL_PASSWORD"), os.Getenv("PSQL_PORT"), os.Getenv("PSQL_DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}

	discBot, err := bot.New(pg)
	if err != nil {
		log.Fatal(err)
	}

	if err = discBot.Start(); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	initModules(router, discBot, pg)

	if err = router.Run(":992"); err != nil {
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("failed starting application: %v", err)))
	}

	<-c
}

func initModules(router *gin.Engine, bot *bot.Bot, db *db.DB) {
	api.InitRoutes(router, services.NewUserService(db))
}
