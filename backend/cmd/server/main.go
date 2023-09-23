package main

import (
	"linked/internal/config"
	"linked/internal/migrations"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Host             string `env:"HOST"`
	AppDBPath        string `env:"APP_DATABASE"`
	AppDBMigrations  string `env:"APP_MIGRATIONS"`
	AuthDBPath       string `env:"AUTH_DATABASE"`
	AuthDBMigrations string `env:"AUTH_MIGRATIONS"`
}

func main() {
	cfg := Config{}

	err := config.LoadEnv(&cfg)
	if err != nil {
		log.Println(err)
		return
	}

	appDB, err := openSqliteDB(cfg.AppDBPath)
	if err != nil {
		log.Println(err)
		return
	}

	err = migrations.RunMigrations(appDB, cfg.AppDBMigrations)
	if err != nil {
		log.Println("error running migrations:", err)
		return
	}
	router := pathrouter.NewRouter()

	router.Group("/api", func(g *pathrouter.Group) {
		g.Use(pathrouter.GzipMiddleware)
		initCollectionsApi(appDB, g)
		initItemsApi(appDB, g)
		initOpenGraphApi(appDB, g)
	})

	log.Println("starting server at", cfg.Host)
	http.ListenAndServe(cfg.Host, router)
}
