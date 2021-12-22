package api

import (
	"github.com/c95rt/context/server"
)

func GetRoutes() []*server.Route {
	return []*server.Route{
		{Path: "/token", Methods: []string{"GET", "HEAD"}, Handler: GenerateToken, IsProtected: false},
		{Path: "/value", Methods: []string{"GET", "HEAD"}, Handler: HandlerWithValue, IsProtected: true},
		{Path: "/cancel", Methods: []string{"GET", "HEAD"}, Handler: HandlerWithCancel, IsProtected: true},
		{Path: "/timeout", Methods: []string{"GET", "HEAD"}, Handler: HandlerWithTimeout, IsProtected: true},
		{Path: "/close", Methods: []string{"GET", "HEAD"}, Handler: HandlerCloseGoRoutine, IsProtected: false},
	}
}
