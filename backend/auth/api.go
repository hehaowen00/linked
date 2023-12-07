package auth

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"image/png"
	"linked/collections"
	"linked/internal/constants"
	"linked/utils"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	pathrouter "github.com/hehaowen00/path-router"
	"github.com/pquerna/otp/totp"
)

type AuthAPI struct {
	authRepo      *AuthRepo
	hmacSecret    string
	collectionApi *collections.CollectionAPI
}

func NewAPI(db *sql.DB, hmacSecret string) *AuthAPI {
	authApi := AuthAPI{
		authRepo:   newAuthRepo(db),
		hmacSecret: hmacSecret,
	}
	return &authApi
}

func (api *AuthAPI) SetCollectionsApi(collectionApi *collections.CollectionAPI) {
	api.collectionApi = collectionApi
}

func (api *AuthAPI) Bind(scope pathrouter.IRoutes) {
	scope.Post("/register", api.Register)
	scope.Post("/login", api.Login)
	scope.Get("/logout", api.Logout)
	scope.Get("/validate", api.Validate)
}

func (api *AuthAPI) Register(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	req := registerRequest{}

	err := utils.ReadJSON(r.Body, &req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  "unable to parse json",
		})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constants.AppName,
		AccountName: req.Email,
	})
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{})
		return
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  "unable to generate qr code",
		})
		return
	}
	png.Encode(&buf, img)

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	userId := uuid.NewString()
	err = api.authRepo.createUser(&req, userId, key.Secret())
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  "failed to create new user",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status":  "ok",
		"qr_code": fmt.Sprintf("data:image/png;base64, %v", encoded),
	})

}

func (api *AuthAPI) Login(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	req := loginRequest{}

	err := utils.ReadJSON(r.Body, &req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  "unable to parse json",
		})
		return
	}

	user, err := api.authRepo.getUser(&req)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  "invalid user details",
		})
		return
	}

	valid := totp.Validate(req.PassCode, user.secret)
	if !valid {
		valid, err = totp.ValidateCustom(req.PassCode, user.secret, time.Now().Add(-30*time.Second), totp.ValidateOpts{})
		if !valid || err != nil {
			log.Println(err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
				"status": "error",
				"error":  "invalid user details",
			})
			return
		}
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	claims := authClaims{
		jwt.RegisteredClaims{
			Subject:   user.id,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString([]byte(api.hmacSecret))
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  "unable to generate access token",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Expires:  expiresAt,
		Secure:   true,
	})

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
	})
}

func (api *AuthAPI) Logout(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	_, err := r.Cookie("access_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.JSON{
			"status": "error",
			"error":  "unauthorized",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		Secure:   true,
	})

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
	})
}

func (api *AuthAPI) Validate(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.JSON{
			"status": "error",
			"error":  "unauthorized",
		})
		return
	}

	token, err := jwt.ParseWithClaims(accessToken.Value, &authClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(api.hmacSecret), nil
	})
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.JSON{
			"status": "error",
			"error":  "unauthorized",
		})
		return
	}

	_, ok := token.Claims.(*authClaims)
	if !token.Valid || !ok {
		return
	}

	// ctx := context.WithValue(r.Context(), constants.AuthKey, claims.Subject)
	// r = r.WithContext(ctx)

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
	})
}

func (api *AuthAPI) Middleware(next pathrouter.HandlerFunc) pathrouter.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		accessToken, err := r.Cookie("access_token")
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.JSON{
				"status": "error",
				"error":  "unauthorized",
			})
			return
		}

		token, err := jwt.ParseWithClaims(accessToken.Value, &authClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(api.hmacSecret), nil
		})
		if err != nil {
			log.Println(err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.JSON{
				"status": "error",
				"error":  "unauthorized",
			})
			return
		}

		claims, ok := token.Claims.(*authClaims)
		if !token.Valid || !ok {
			return
		}

		ctx := context.WithValue(r.Context(), constants.AuthKey, claims.Subject)
		r = r.WithContext(ctx)

		next(w, r, ps)
	}
}
