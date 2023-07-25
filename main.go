package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MixinNetwork/safe/governance/blaze"
	"github.com/MixinNetwork/safe/governance/cmd"
	"github.com/MixinNetwork/safe/governance/config"
	"github.com/MixinNetwork/safe/governance/middlewares"
	"github.com/MixinNetwork/safe/governance/routes"
	"github.com/MixinNetwork/safe/governance/session"
	"github.com/MixinNetwork/safe/governance/store"
	"github.com/dimfeld/httptreemux"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "governance",
		Usage:                "Mixin Safe Governance",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:   "http",
				Usage:  "Run the http service",
				Action: bootCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "environment",
						Aliases: []string{"e"},
						Value:   "development",
						Usage:   "The environment of the http service",
					},
				},
			},
			{
				Name:   "migrate",
				Usage:  "Migrate the app's ownership",
				Action: cmd.MigrateCMD,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "keystore",
						Aliases: []string{"k"},
						Usage:   "The encrypted keystore of the app",
					},
					&cli.StringFlag{
						Name:    "private",
						Aliases: []string{"s"},
						Usage:   "The private spend key of the custodian",
					},
					&cli.StringFlag{
						Name:    "public",
						Aliases: []string{"p"},
						Usage:   "The public key of the app, can be found from https://governance.mixin.one, for staging is https://governance.mixin.zone",
					},
					&cli.StringFlag{
						Name:    "user",
						Aliases: []string{"u"},
						Usage:   "The user who will receive the app",
					},
					&cli.StringFlag{
						Name:    "encrypted",
						Aliases: []string{"e"},
						Value:   "true",
						Usage:   "The keystore is encrypted",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func bootCmd(c *cli.Context) error {
	env := c.String("environment")
	config.InitConfiguration(env)

	database, err := store.OpenDatabase()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx = session.WithDatabase(ctx, database)
	go blaze.Boot(ctx)

	router := httptreemux.New()
	routes.RegisterRoutes(router)
	handler := middlewares.Constraint(router)
	handler = middlewares.Context(handler, database, render.New())
	handler = middlewares.Stats(handler)

	log.Printf("Mixin Safe Governance http service start at: http://localhost:%s\n", config.AppConfig.Port)
	return http.ListenAndServe(fmt.Sprintf(":%s", config.AppConfig.Port), handler)
}
