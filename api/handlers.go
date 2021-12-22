package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/c95rt/context/config"
	"github.com/c95rt/context/helpers"
	"github.com/c95rt/context/models"
	"github.com/c95rt/context/server"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

func GenerateToken(ctx *config.AppContext, w *server.ResponseWriter, r *http.Request) {
	claims := struct {
		User map[string]interface{} `json:"u"`
		jwt.StandardClaims
	}{
		map[string]interface{}{
			"r":     []int{1},
			"email": "test.user@globant.com",
		},
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}
	jwtSecret := "asdf"
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	if err != nil {
		w.WriteJSON(http.StatusInternalServerError, err.Error())
	}

	w.WriteJSON(http.StatusOK, token)
}

func HandlerWithValue(ctx *config.AppContext, w *server.ResponseWriter, r *http.Request) {
	funcName := "HandlerWithValue"
	userInfo := models.InfoUser{}

	requestContext := r.Context()
	contextValue := requestContext.Value("user")
	if err := mapstructure.Decode(contextValue, &userInfo); err != nil {
		fmt.Printf("%s: error at mapstructure.Decode: %s\n", funcName, err.Error())
		w.WriteJSON(http.StatusInternalServerError, err)
	}

	w.WriteJSON(http.StatusOK, userInfo)
}

func HandlerWithCancel(ctx *config.AppContext, w *server.ResponseWriter, r *http.Request) {
	funcName := "HandlerWithCancel"

	requestContext := r.Context()

	requestContext, cancel := context.WithCancel(requestContext)

	response := &models.TurnResponse{}

	go func(parentContext context.Context) {
		for {
			closeGoRoutines := <-ctx.CloseGoRoutine
			if closeGoRoutines {
				fmt.Printf("%s: go routine was closed by another handler\n", funcName)
				cancel()
				break
			}
		}
	}(requestContext)

	helpers.PrintAfterTime(requestContext, response)

	w.WriteJSON(http.StatusOK, response)
}

func HandlerCloseGoRoutine(ctx *config.AppContext, w *server.ResponseWriter, r *http.Request) {
	handlerName := "HandlerCloseGoRoutine"
	fmt.Printf("%s: go routine will be closed\n", handlerName)
	go func() {
		ctx.CloseGoRoutine <- true
	}()
	fmt.Printf("%s: go routine was closed\n", handlerName)
	return
}

func HandlerWithTimeout(ctx *config.AppContext, w *server.ResponseWriter, r *http.Request) {
	requestContext := r.Context()

	var duration int = 3

	requestContext, cancel := context.WithTimeout(requestContext, time.Duration(duration)*time.Second)
	defer cancel()

	response := &models.TurnResponse{}
	helpers.PrintAfterTime(requestContext, response)

	w.WriteJSON(http.StatusOK, response)
}
