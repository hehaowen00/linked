package main

import (
	"context"
	"database/sql"
	"io"
	"linked/internal/constants"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	pathrouter "github.com/hehaowen00/path-router"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func initAuthApi(db *sql.DB, auth *GoogleAuth, scope pathrouter.IRoutes) {
	scope.Get("/login", auth.login)
	scope.Get("/callback", auth.handleCallback)
	scope.Get("/validate", auth.ValidateToken)
	scope.Get("/logout", auth.logout)
	scope.Get("/profile", auth.getProfile)
}

type GoogleAuth struct {
	config       *oauth2.Config
	csrf         string
	db           *sql.DB
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
		defer r.Body.Close()

		accessToken, err := r.Cookie(constants.AccessTokenKey)
		if err != nil {
			log.Println("no access token", err)
			http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
			return
		}

		response, err := http.Get(
			constants.GoogleOAuthTokenInfoUrl + accessToken.Value,
		)
		if err != nil {
			log.Println("token info", err)
			http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
			return
		}

		if response.StatusCode != 200 {
			http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
			return
		}

		contents, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
			return
		}

		info := TokenInfo{}

		err = json.Unmarshal(contents, &info)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
			return
		}

		if info.Audience != auth.config.ClientID {
			log.Println("incorrect audience")
			http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
			return
		}

		c := setContext(r.Context(), constants.AuthKey, info.Subject)
		r = r.WithContext(c)
		next(w, r, ps)
	}
}

func (auth *GoogleAuth) login(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	defer r.Body.Close()

	redirect := r.URL.Query().Get(constants.RedirectUrlKey)
	state := auth.csrf
	if redirect != "" {
		state = state + "&" + redirect
	}
	authUrl, err := url.Parse(auth.config.AuthCodeURL(state))
	if err != nil {
		log.Println(err)
	}
	authUrl.Query().Add("approval_prompt", "auto")
	url := auth.config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (auth *GoogleAuth) handleCallback(
	w http.ResponseWriter,
	r *http.Request,
	ps *pathrouter.Params,
) {
	defer r.Body.Close()

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

	http.SetCookie(w, &http.Cookie{
		Name:    constants.AccessTokenKey,
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
	defer r.Body.Close()

	accessToken, err := r.Cookie(constants.AccessTokenKey)
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
	defer r.Body.Close()

	accessToken, err := r.Cookie(constants.AccessTokenKey)
	if err != nil {
		log.Println("no access token", err, accessToken)
		writeJson(w, http.StatusUnauthorized, JsonResult{
			Status: "error",
			Error:  "missing access token",
		})
		return
	}
	response, err := http.Get(
		constants.GoogleOAuthTokenInfoUrl + accessToken.Value,
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
		Data:   tokenInfo,
	})
}

func (auth *GoogleAuth) logout(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	defer r.Body.Close()

	http.SetCookie(w, &http.Cookie{
		Name:    constants.AccessTokenKey,
		Value:   "",
		Expires: time.Now().UTC().Add(-time.Minute),
		Path:    "/",
	})
	http.Redirect(w, r, auth.frontendHost, http.StatusTemporaryRedirect)
}

func (auth *GoogleAuth) getProfile(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	defer r.Body.Close()

	info, err := getUserInfo(r)
	if err != nil {
	}

	writeJson(w, http.StatusOK, info)
}

func getUserInfo(r *http.Request) (*UserInfo, error) {
	accessToken, err := r.Cookie(constants.AccessTokenKey)
	if err != nil {
		log.Println("no access token", err)
		return nil, err
	}
	response, err := http.Get(
		constants.GoogleOAuthUserInfoUrl + accessToken.Value,
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
