package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	pathrouter "github.com/hehaowen00/path-router"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func initAuthApi(db *sql.DB, auth *GoogleAuth, router *pathrouter.Group) {
	router.Get("/login", auth.login)
	router.Get("/callback", auth.handleCallback)
	router.Get("/validate", auth.ValidateToken)
	router.Get("/logout", auth.logout)
	router.Get("/profile", auth.getProfile)
}

type GoogleAuth struct {
	config       *oauth2.Config
	csrf         string
	db           *sql.DB
	authHost     string
	frontendHost string
}

func newGoogleAuth(host, clientId, clientSecret string, db *sql.DB) *GoogleAuth {
	return &GoogleAuth{
		config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  host + "/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		csrf: uuid.NewString(),
		db:   db,
	}
}

func (auth *GoogleAuth) authMiddleware(next pathrouter.HandlerFunc) pathrouter.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		accessToken, err := r.Cookie("access_token")
		if err != nil {
			log.Println("no access token", err)
			writeJson(w, http.StatusUnauthorized, JsonResult{
				Status: "error",
				Error:  "missing access token",
			})
			return
		}

		response, err := http.Get(
			"https://oauth2.googleapis.com/tokeninfo?access_token=" + accessToken.Value,
		)
		if err != nil {
			log.Println("token info", err)
			return
		}

		if response.StatusCode != 200 {
			return
		}

		contents, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			return
		}

		info := TokenInfo{}

		err = json.Unmarshal(contents, &info)
		if err != nil {
			log.Println(err)
		}

		if info.Audience != auth.config.ClientID {
			log.Println("incorrect audience")
			return
		}

		c := setContext(r.Context(), "id", info.Subject)
		r = r.WithContext(c)
		next(w, r, ps)
	}
}

func (auth *GoogleAuth) login(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	redirect := r.URL.Query().Get("redirect_url")
	url := auth.config.AuthCodeURL(auth.csrf + "&" + redirect)
	http.Redirect(
		w,
		r,
		url+"&approval_prompt=auto",
		http.StatusTemporaryRedirect,
	)
}

func (auth *GoogleAuth) handleCallback(
	w http.ResponseWriter,
	r *http.Request,
	ps *pathrouter.Params,
) {
	stateValue := r.FormValue("state")
	codeValue := r.FormValue("code")

	if !strings.HasPrefix(stateValue, auth.csrf) {
		log.Println("oauth2 flow - CSRF failed")
		writeJson(w, http.StatusUnauthorized, JsonResult{
			Status: "error",
			Error:  "CSRF check failed",
		})
		return
	}

	redirectUrl := stateValue[len(auth.csrf):]

	token, err := auth.config.Exchange(context.Background(), codeValue)
	if err != nil {
		log.Println("code exchanged failed:", err)
		writeJson(w, http.StatusUnauthorized, JsonResult{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	response, err := http.Get(
		"https://oauth2.googleapis.com/tokeninfo?access_token=" + token.AccessToken,
	)
	if err != nil {
		log.Println(err)
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	tokenInfo := &TokenInfo{}
	err = json.Unmarshal(contents, &tokenInfo)
	if err != nil {
		log.Println(err)
		return
	}

	if tokenInfo.Audience != auth.config.ClientID {
		writeJson(w, http.StatusUnauthorized, JsonResult{
			Status: "error",
			Error:  "unauthorized",
		})
		return
	}

	// http.SetCookie(w, &http.Cookie{
	// 	Name:  "refresh_token",
	// 	Value: token.RefreshToken,
	// 	Path:  "/",
	// })

	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   token.AccessToken,
		Expires: token.Expiry,
		Path:    "/",
	})

	if redirectUrl != "" {
		http.Redirect(w, r, redirectUrl[1:], http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
}

func (auth *GoogleAuth) getUserInfo(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		log.Println("no access token", err, accessToken)
		writeJson(w, http.StatusUnauthorized, JsonResult{
			Status: "error",
			Error:  "missing access token",
		})
		return
	}
	response, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken.Value,
	)
	if err != nil {
		log.Println(err)
		writeJson(w, http.StatusUnauthorized, JsonResult{
			Status: "error",
			Error:  "failed to get user info",
		})
		return
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	userInfo := UserInfo{}
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		log.Println(err)
		return
	}

	writeJson(w, http.StatusOK, userInfo)
}

func (auth *GoogleAuth) ValidateToken(
	w http.ResponseWriter,
	r *http.Request,
	ps *pathrouter.Params,
) {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		log.Println("no access token", err, accessToken)
		writeJson(w, http.StatusUnauthorized, JsonResult{
			Status: "error",
			Error:  "missing access token",
		})
		return
	}
	response, err := http.Get(
		"https://oauth2.googleapis.com/tokeninfo?access_token=" + accessToken.Value,
	)
	if err != nil {
		log.Println(err)
		writeJson(w, http.StatusInternalServerError, JsonResult{
			Status: "error",
		})
		return
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		writeJson(w, http.StatusInternalServerError, JsonResult{
			Status: "error",
		})
		return
	}

	tokenInfo := TokenInfo{}
	err = json.Unmarshal(contents, &tokenInfo)
	if err != nil {
		log.Println(err)
		writeJson(w, http.StatusInternalServerError, JsonResult{
			Status: "error",
		})
		return
	}

	if tokenInfo.Audience != auth.config.ClientID {
		log.Println("audience does not match client id")
		writeJson(w, http.StatusUnauthorized, JsonResult{
			Status: "error",
			Error:  "unauthorized",
		})
		return
	}

	writeJson(w, http.StatusOK, JsonResult{
		Status: "ok",
	})
}

func (auth *GoogleAuth) logout(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: time.Now().UTC().Add(-time.Minute),
		Path:    "/",
	})
	http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
}

func (auth *GoogleAuth) getProfile(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	info, err := getUserInfo(r)
	if err != nil {
	}

	writeJson(w, http.StatusOK, info)
}

func getUserInfo(r *http.Request) (*UserInfo, error) {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		log.Println("no access token", err)
		return nil, err
	}
	response, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken.Value,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	userInfo := UserInfo{}
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &userInfo, nil
}
