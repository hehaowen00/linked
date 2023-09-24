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
	Host string `env:"HOST"`

	AppDBPath       string `env:"APP_DATABASE"`
	AppDBMigrations string `env:"APP_MIGRATIONS"`

	AuthDBPath       string `env:"AUTH_DATABASE"`
	AuthDBMigrations string `env:"AUTH_MIGRATIONS"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`

	AuthHost     string `env:"AUTH_HOST"`
	FrontendHost string `env:"FRONTEND_HOST"`
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

	authDB, err := openSqliteDB(cfg.AuthDBPath)
	if err != nil {
		log.Println(err)
		return
	}

	err = migrations.RunMigrations(authDB, cfg.AuthDBMigrations)
	if err != nil {
		log.Println("error running migrations:", err)
		return
	}

	googleAuth := newGoogleAuth(cfg.AuthHost, cfg.GoogleClientID, cfg.GoogleClientSecret, authDB)
	googleAuth.frontendHost = cfg.FrontendHost

	router := pathrouter.NewRouter()

	router.Group("/auth", func(g *pathrouter.Group) {
		g.Use(Cors)
		initAuthApi(authDB, googleAuth, g)
	})

	router.Group("/api", func(g *pathrouter.Group) {
		g.Use(pathrouter.GzipMiddleware, googleAuth.authMiddleware, Cors)
		initCollectionsApi(appDB, g)
		initItemsApi(appDB, g)
		initOpenGraphApi(appDB, g)
	})

	router.Handle(
		http.MethodOptions,
		"*",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		}),
	)

	log.Println("starting server at", cfg.Host)
	http.ListenAndServe(cfg.Host, router)
}

func Cors(next pathrouter.HandlerFunc) pathrouter.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next(w, r, ps)
	}
}
