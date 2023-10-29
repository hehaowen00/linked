package main

import (
	"io"
	"linked/auth"
	"linked/collections"
	"linked/internal/config"
	"linked/items"
	"linked/opengraph"
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
	AppName   string `env:"APP_NAME"`
	Host      string `env:"HOST"`
	StaticDir string `env:"STATIC_DIR"`

	AppDBPath       string `env:"APP_DATABASE"`
	AppDBMigrations string `env:"APP_MIGRATIONS"`

	AuthDBPath       string `env:"AUTH_DATABASE"`
	AuthDBMigrations string `env:"AUTH_MIGRATIONS"`
	AuthSecret       string `env:"AUTH_SECRET"`
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

	router := pathrouter.NewRouter()
	router.Use(pathrouter.GzipMiddleware)

	var cors pathrouter.CorsHandler
	cors.AllowCredentials = true
	router.Use(cors.Middleware)

	router.Get("/app", serveIndexHtml(cfg.StaticDir))
	router.Get("/app/*", serveIndexHtml(cfg.StaticDir))
	router.Get("/*", createSPAHandler(cfg.StaticDir))

	authScope := router.Scope("/auth")

	authApi := auth.NewAPI(authDB, "secret")
	authApi.Bind(authScope)

	apiScope := router.Scope("/api")
	apiScope.Use(authApi.Middleware)

	collectionsApi := collections.NewAPI(appDB)
	collectionsApi.Bind(apiScope)
	authApi.SetCollectionsApi(collectionsApi)

	itemsApi := items.NewAPI(appDB)
	itemsApi.Bind(apiScope)

	opengraph.InitOpenGraphApi(appDB, apiScope)

	log.Println("starting server at", cfg.Host)
	go func() {
		http.HandleFunc("/", redirectToHTTPS)
		log.Fatalln(http.ListenAndServe(":80", nil))
	}()
	log.Fatalln(http.ListenAndServeTLS(cfg.Host, "server.crt", "server.key", router))
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), http.StatusFound)
		return
	}
}

func serveIndexHtml(rootDir string) pathrouter.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		fp := filepath.Join(rootDir, "index.html")
		f, err := os.OpenFile(fp, os.O_RDONLY, 0777)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext("index.html"))
		w.Header().Add("Content-Type", contentType)
		io.Copy(w, f)
	}
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
