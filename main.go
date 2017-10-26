package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/orendo/gin-tokenauth"
	"github.com/orendo/gin-tokenauth/filestore"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "image-admission"
	app.Version = "v0.0.1"
	app.Usage = "container image admission service"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "enable db query logging",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "listen",
			Usage:  "listen ip:port",
			EnvVar: "LISTEN",
			Value:  ":8000",
		},
		cli.StringFlag{
			Name:   "tokens-file",
			Usage:  "path to tokens file for Bearer token-based authentication",
			EnvVar: "TOKENS_FILE",
		},
		cli.StringFlag{
			Name:   "dbhost",
			Usage:  "database host",
			EnvVar: "DBHOST",
			Value:  "localhost",
		},
		cli.IntFlag{
			Name:   "dbport",
			Usage:  "database port",
			EnvVar: "DBPORT",
			Value:  5432,
		},
		cli.StringFlag{
			Name:   "dbuser",
			Usage:  "database user",
			EnvVar: "DBUSER",
			Value:  "postgres",
		},
		cli.StringFlag{
			Name:   "dbpassword",
			Usage:  "database password",
			EnvVar: "DBPASSWORD",
			Value:  "postgres",
		},
		cli.StringFlag{
			Name:   "dbname",
			Usage:  "database name",
			EnvVar: "DBNAME",
			Value:  "imageadmission",
		},
		cli.StringFlag{
			Name:   "dbsslmode",
			Usage:  "database sslmode",
			EnvVar: "DBSSLMODE",
			Value:  "disable",
		},
		cli.StringFlag{
			Name:   "certfile",
			Usage:  "TLS certificate file name",
			EnvVar: "CERTFILE",
		},
		cli.StringFlag{
			Name:   "keyfile",
			Usage:  "TLS certificate private key file name",
			EnvVar: "KEYFILE",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		dbConnString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
			ctx.String("dbuser"),
			ctx.String("dbpassword"),
			ctx.String("dbname"),
			ctx.String("dbsslmode"),
			ctx.String("dbhost"),
			strconv.Itoa(ctx.Int("dbport")),
		)

		db, err := gorm.Open("postgres", dbConnString)
		if err != nil {
			return err
		}
		defer db.Close()

		db.LogMode(ctx.Bool("debug"))

		if err := db.AutoMigrate(&Image{}).Error; err != nil {
			return err
		}

		r := gin.Default()

		r.GET("/images", getImages(db))
		r.GET("/images/:id", getImages(db))

		authorized := r.Group("/")
		if f := ctx.String("tokens-file"); f != "" {
			store, err := filestore.New(f)
			if err != nil {
				log.Fatalln(err)
			}
			authorized.Use(tokenauth.New(store))
		}

		authorized.PUT("/images", putImage(db))
		authorized.DELETE("/images/:id", deleteImage(db))

		if ctx.String("certfile") != "" && ctx.String("keyfile") != "" {
			return r.RunTLS(ctx.String("listen"), ctx.String("certfile"), ctx.String("keyfile"))
		}
		return r.Run(ctx.String("listen"))
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
