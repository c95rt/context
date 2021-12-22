package main

import (
	"log"
	"os"

	"github.com/c95rt/context/api"
	"github.com/c95rt/context/server"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {
	_ = godotenv.Load("prod.env")

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "api-up",
			Usage: "This command starts the service",
			Action: func(c *cli.Context) error {
				server.UpServer(api.GetRoutes())
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
