package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"milton/pkg/db"
	"milton/pkg/db/orm"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/urfave/cli/v2"
)

const (
	dbURL = "./milton.database"
)

var (
	NodeList = []string{"og:site_name"}
)

func Start() {
	app := &cli.App{
		Name:  "milton",
		Usage: "SEO Scraper",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    "list",
				Aliases: []string{"l"},
			},
		},
		Action: func(ctx *cli.Context) error {
			if len(ctx.String("url")) > 0 {
				return Crawl(ctx.Context, ctx.String("url"))
			}

			if len(ctx.String("list")) > 0 {
				list, err := ioutil.ReadFile(ctx.String("list"))
				if err != nil {
					return err
				}

				files := strings.Split(string(list), "\n")

				for _, f := range files {
					err = Crawl(ctx.Context, f)
					if err != nil {
						log.Printf("ERROR: Failed to crawl %s, %v\n", f, err)
					}
				}
			}

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func prettyJson(v interface{}) {
	b, _ := json.MarshalIndent(v, "\t", "\t")
	fmt.Println(string(b))
}

func Crawl(ctx context.Context, url string) error {
	db := db.New(&db.SqliteOpts{ConnString: dbURL})

	// --- Local file debug
	// dir, err := filepath.Abs(filepath.Dir("."))
	// if err != nil {
	// 	panic(err)
	// }

	c := colly.NewCollector()

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c.WithTransport(t)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML("head", func(h *colly.HTMLElement) { // err := db.WithinTransaction(ctx.Context, func(ctx context.Context) error {
		err := db.WithinTransaction(ctx, func(txCtx context.Context) error {
			ws, txErr := db.WithQtx(txCtx).UpsertWebsite(txCtx, h.Request.URL.String())
			if txErr != nil {
				fmt.Println(txErr.Error())
				return txErr
			}

			for _, v := range h.DOM.Children().Nodes {
				if v.Data == "script" {
					continue
				}

				t := orm.UpsertTagParams{WebsiteID: ws.WebsiteID}
				t.TagType = sql.NullString{String: v.Data, Valid: true}
				for _, v := range v.Attr {
					if v.Key == "property" {
						t.Property = sql.NullString{String: v.Val, Valid: true}
					}

					if v.Key == "name" {
						t.Name = sql.NullString{String: v.Key, Valid: true}
					}

					if v.Key == "content" {
						t.Value = sql.NullString{String: v.Val, Valid: true}
					}
				}
				_, txErr = db.WithQtx(txCtx).UpsertTag(txCtx, t)
				if txErr != nil {
					fmt.Println(txErr.Error())
					return txErr
				}
			}

			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("ERROR!", err)
	})

	// c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.Visit(url)
	// ---- URL
	// c.Visit("file://" + dir + "/index.html")
	return nil
}
