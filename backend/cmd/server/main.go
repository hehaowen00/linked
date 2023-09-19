package main

import (
	"database/sql"
	"linked/internal/migrations"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := make(map[string]string)

	err := LoadEnvConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", config["DATABASE_PATH"])
	if err != nil {
		return
	}

	err = migrations.RunMigrations(db, config["MIGRATIONS_PATH"])
	if err != nil {
		log.Println("error running migrations:", err)
		return
	}

	router := pathrouter.NewRouter()

	router.Group("/api", func(g *pathrouter.Group) {
		g.Use(pathrouter.GzipMiddleware)
		initCollectionsApi(db, g)
		initItemsApi(db, g)
	})

	log.Println("starting server at", config["ADDR"])
	http.ListenAndServe(config["ADDR"], router)
}
