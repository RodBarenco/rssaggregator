package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"time"

	database "github.com/RodBarenco/rssaggregator/db"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type APIapiConfig struct {
	DB *gorm.DB
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing environment file path argument")
	}

	envFilePath := os.Args[1]

	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Erro ao carregar variáveis de ambiente:", err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not found in the environment")
	}

	db, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&database.User{},
		&database.Feed{},
		&database.FeedFollows{},
		&database.Post{},
	)

	if err != nil {
		log.Fatal("Failed to apply migrations:", err)
	}

	apiCfg := APIapiConfig{
		DB: db,
	}

	go startScraping(&apiCfg, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsFromUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedsFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedsFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	if envFilePath == ".test.env" {
		fmt.Printf("VERSÃO DE TESTE. ")
	} else {
		fmt.Printf("\033[33mVERSÃO DE DESENVOLVIMENTO. \033[0m")
	}
	log.Printf("\033[33mServer starting on PORT: %v\033[0m", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
