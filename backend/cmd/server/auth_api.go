package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	pathrouter "github.com/hehaowen00/path-router"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type TokenInfo struct {
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
}

type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail string `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type StateInfo struct {
	CSRFToken   string `csv:"csrf"`
	RedirectUrl string `csv:"redirect_url"`
}

func initAuthApi(db *sql.DB, auth *GoogleAuth, router *pathrouter.Group) {
	router.Get("/login", auth.login)
	router.Get("/callback", auth.handleCallback)
	router.Get("/validate", auth.ValidateToken)
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
			log.Println("token info", err)
			return
		}

		contents, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(contents))

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
		r.WithContext(c)
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

	log.Println("stateValue", stateValue)
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

	fmt.Println(string(contents))

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
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	tokenInfo := TokenInfo{}
	err = json.Unmarshal(contents, &tokenInfo)
	if err != nil {
	}

	if tokenInfo.Audience != auth.config.ClientID {
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

func validateToken(r *http.Request) error {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		log.Println("no access token", err, accessToken)
		return err
	}
	response, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken.Value,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	userInfo := UserInfo{}
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
