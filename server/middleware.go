package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/c95rt/context/config"
	"github.com/c95rt/context/helpers"
	"github.com/c95rt/context/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/lithammer/shortuuid/v3"
	jwtmiddleware "github.com/mfuentesg/go-jwtmiddleware"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	ROLE_ADMIN int = 1
	ROLE_USER  int = 2
)

func jwtErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	r := &ResponseWriter{writer: w}
	if err.Error() == "Token is expired" {
		r.Error(http.StatusUnauthorized, "unauthorized", WithErrorScope("token"), WithErrorType(1))
		return
	}
	if err != nil {
		r.Error(http.StatusUnauthorized, "unauthorized", WithErrorScope("token"))
	}
}

func NewJWTMiddleware(secret []byte) *jwtmiddleware.Middleware {
	return jwtmiddleware.New(
		jwtmiddleware.WithErrorHandler(jwtErrorHandler),
		jwtmiddleware.WithSigningMethod(jwt.SigningMethodHS256),
		jwtmiddleware.WithSignKey(secret),
		jwtmiddleware.WithUserProperty("_jwt-token"),
	)
}

func LoggerRequest(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestLogger := log.WithFields(
		log.Fields{
			"request_id":         r.Header.Get("X-Request-ID"),
			"request_id_service": shortuuid.New(),
			"query":              r.URL.Query(),
			"host":               r.Host,
			"url":                r.URL.Path,
			"headers":            r.Header})
	requestLogger.Info("logger_request")
	config.SetLogger(requestLogger)
	next(rw, r)
}

func UserMiddleware() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		authorization := r.Header.Get("Authorization")
		token := strings.Split(authorization, " ")
		if len(token) == 2 {
			tokenString := token[1]
			data, _ := helpers.ParserTokenUnverified(tokenString)
			tokenParse, ok := data["u"].(map[string]interface{})
			if ok {
				dataInfo := models.InfoUser{}
				roles := tokenParse["r"]
				mapstructure.Decode(map[string]interface{}{
					"Roles": roles,
				}, &dataInfo)
				isAdmin := helpers.Contains(dataInfo.Roles, ROLE_ADMIN)
				isUser := helpers.Contains(dataInfo.Roles, ROLE_USER)
				userInfo := map[string]interface{}{
					"Email":   tokenParse["email"],
					"IsAdmin": isAdmin,
					"IsUser":  isUser,
					"Roles":   roles,
				}
				if !isAdmin && !isUser {
					a := &ResponseWriter{writer: rw}
					a.Error(http.StatusUnauthorized, "unauthorized", WithErrorScope("token"))
					return
				}
				ctx := context.WithValue(r.Context(), string("user"), userInfo)
				next(rw, r.WithContext(ctx))
				return
			}
		}
		next(rw, r)
	})
}
