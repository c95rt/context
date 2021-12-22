package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/c95rt/context/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

func ParserTokenUnverified(tokenStr string) (jwt.MapClaims, bool) {
	var p jwt.Parser
	token, _, ok := p.ParseUnverified(tokenStr, jwt.MapClaims{})
	if ok != nil {
		return nil, false
	}
	tokendata, _ := token.Claims.(jwt.MapClaims)
	return tokendata, true
}

func Contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func PrintAfterTime(parentContext context.Context, response *models.TurnResponse) {
	funcName := "PrintAfterTime"
	userInfo := models.InfoUser{}
	contextValue := parentContext.Value("user")
	if err := mapstructure.Decode(contextValue, &userInfo); err != nil {
		fmt.Printf("%s: error at mapstructure.Decode: %s\n", funcName, err.Error())
		return
	}
	for {
		select {
		case <-parentContext.Done():
			fmt.Println("PrintAfterTime finished")
			return
		default:
			response.Turn += 1
			fmt.Printf("%s: %s turn: %d\n", funcName, userInfo.Email, response.Turn)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
