package main

import (
	"fmt"
	"main/tools"
	"os"

	"github.com/restream/reindexer"
	"github.com/urfave/cli/v2"
)

type Item struct {
	ID   int64  `reindex:"id,,pk"`
	Name string `reindex:"name"`
	Year int    `reindex:"year,tree"`
}

func main() {
	config, err := tools.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db := reindexer.NewReindex(config.ServerAddress + "/" + config.DatabaseName)

	db.OpenNamespace(config.Namespace, reindexer.DefaultNamespaceOptions(), Item{})

	app := &cli.App{
		Name:  "reindex-db-loader",
		Usage: "Reindex Database loader",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add a document",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
					&cli.IntFlag{Name: "year", Aliases: []string{"y"}},
				},
				Action: func(cCtx *cli.Context) error {
					// Find the ID of the new document
					query := db.Query(config.Namespace).ReqTotal()
					iterator := query.Exec()
					defer iterator.Close()

					if err := iterator.Error(); err != nil {
						panic(err)
					}

					LastID := int64(iterator.TotalCount())

					// Add a document
					err := db.Upsert(config.Namespace, &Item{
						ID:   LastID,
						Name: cCtx.String("name"),
						Year: cCtx.Int("year"),
					})
					if err != nil {
						panic(err)
					}

					// Query of the added document
					elem, found := db.Query(config.Namespace).
						Where("id", reindexer.EQ, LastID).
						Get()

					if found {
						item := elem.(*Item)
						fmt.Println("Added document:", *item)
					}

					return nil
				},
			},
			{
				Name:    "getlist",
				Aliases: []string{"gl"},
				Usage:   "Get a list of documents",
				Action: func(cCtx *cli.Context) error {
					query := db.Query(config.Namespace).
						Sort("id", false).
						ReqTotal()

					iterator := query.Exec()
					defer iterator.Close()

					fmt.Println("Found", iterator.TotalCount(), "total documents. Documents:")

					// Iterate over results
					for iterator.Next() {
						elem := iterator.Object().(*Item)
						fmt.Println(*elem)
					}

					if err := iterator.Error(); err != nil {
						panic(err)
					}

					return nil
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Get the document",
				Flags: []cli.Flag{
					&cli.Int64Flag{Name: "id", Aliases: []string{"i"}},
				},
				Action: func(cCtx *cli.Context) error {
					elem, found := db.Query(config.Namespace).
						Where("id", reindexer.EQ, cCtx.Int64("id")).
						Get()

					if found {
						item := elem.(*Item)
						fmt.Println("Document:", *item)
					} else {
						fmt.Println("There is no such document.")
					}

					return nil
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update document",
				Flags: []cli.Flag{
					&cli.Int64Flag{Name: "id", Aliases: []string{"i"}},
					&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
					&cli.IntFlag{Name: "year", Aliases: []string{"y"}},
				},
				Action: func(cCtx *cli.Context) error {
					// Update document
					err := db.Upsert(config.Namespace, &Item{
						ID:   cCtx.Int64("id"),
						Name: cCtx.String("name"),
						Year: cCtx.Int("year"),
					})
					if err != nil {
						panic(err)
					}

					// Request an updated document
					elem, found := db.Query(config.Namespace).
						Where("id", reindexer.EQ, cCtx.Int64("id")).
						Get()

					if found {
						item := elem.(*Item)
						fmt.Println("Updated document:", *item)
					}

					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Delete a document",
				Flags: []cli.Flag{
					&cli.Int64Flag{Name: "id", Aliases: []string{"i"}},
				},
				Action: func(cCtx *cli.Context) error {
					err := db.Delete(config.Namespace, &Item{
						ID: cCtx.Int64("id"),
					})
					if err != nil {
						panic(err)
					}

					fmt.Println("delete task:", cCtx.String("id"))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
