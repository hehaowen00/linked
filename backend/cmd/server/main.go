package main

import (
	"io"
	"linked/internal/config"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	pathrouter "github.com/hehaowen00/path-router"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Host      string `env:"HOST"`
	StaticDir string `env:"STATIC_DIR"`

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

	appDB, err := loadDB(cfg.AppDBPath, cfg.AppDBMigrations)
	if err != nil {
		log.Println(err)
		return
	}

	authDB, err := loadDB(cfg.AuthDBPath, cfg.AuthDBMigrations)
	if err != nil {
		log.Println("error running migrations:", err)
		return
	}

	googleAuth := newGoogleAuth(cfg.AuthHost, cfg.GoogleClientID, cfg.GoogleClientSecret, authDB)
	googleAuth.frontendHost = cfg.FrontendHost

	router := pathrouter.NewRouter()
	router.Use(pathrouter.GzipMiddleware)
	router.Get("/*", createSPAHandler(cfg.StaticDir))

	authScope := router.Scope("/auth")

	var cors pathrouter.CorsHandler
	cors.AllowCredentials = true
	authScope.Use(cors.Middleware)

	initAuthApi(authDB, googleAuth, authScope)

	apiScope := router.Scope("/api")
	apiScope.Use(googleAuth.authMiddleware)

	initCollectionsApi(appDB, apiScope)
	initItemsApi(appDB, apiScope)
	initOpenGraphApi(appDB, apiScope)

	log.Println("starting server at", cfg.Host)
	log.Fatalln(http.ListenAndServe(cfg.Host, router))
}

func createSPAHandler(rootDir string) pathrouter.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		path := ps.Get("*")

		if strings.Contains(path, "..") {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if len(path) == 0 {
			path = "index.html"
		}

		fp := filepath.Join(rootDir, path)
		f, err := os.OpenFile(fp, os.O_RDONLY, 0777)
		if err != nil {
			fp := filepath.Join(rootDir, "index.html")
			f, err = os.OpenFile(fp, os.O_RDONLY, 0777)
		}

		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(path))
		w.Header().Add("Content-Type", contentType)
		io.Copy(w, f)
	}
}
