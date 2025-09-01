package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	routes "github.com/suhas-developer07/WebscraperAndContenanalysis/internal"
	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/database"
	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/repository"
)

// --------Project Blueprint-------

// Phase-1 (should be complete within 30-8 to 2-9)

// 1. Create an endpoint that accepts urls like this {"urls": ["https://example.com/1", "https://example.com/2"]}.
// 2. Connect To Postgress and insert the incoming urls into DB
// 3. Connect to RabbitMQ to publish message for each task to the queue. the message should be looks like this id and urls.

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DatabaseUrl := os.Getenv("DatabaseUrl")

	if DatabaseUrl == "" {
		log.Fatal("DatabaseUrl not found in env")
	}

	Database, err := database.Connect(DatabaseUrl)

	if err != nil {
		log.Fatal("failed to initialize database: %w", err)
	}
	defer Database.Close()

	repo := repository.NewPostgresRepository(Database)

	if err := repo.InitTable(); err != nil {
		log.Fatalln("Error initializing the table", err)
	}

	routes := routes.MountRoutes(*repo)

	http.ListenAndServe(":8080", routes)

}
