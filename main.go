// ЗАПУСТИТЬ DOCKER
// go build . && main
package main

import (
	"fmt"
	"main/tools"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/restream/reindexer"
)

type Item struct {
	ID   int64  `reindex:"id,,pk"`
	Name string `reindex:"name"`
	Year int    `reindex:"year,tree"`
}

func main() {
	router := gin.Default()

	config, err := tools.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db := reindexer.NewReindex(config.ReindexerServerAddress + "/" + config.DatabaseName)

	db.OpenNamespace(config.Namespace, reindexer.DefaultNamespaceOptions(), Item{})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, `COMMANDS:
		/		Shows a list of commands or help for one command
		/add/:name/:year		Add a document
		/get/:id		Get the document
		/getlist		Get a list of documents
		/update/:id/:name/:year		Update document
		/delete/:id		Delete a document`)
	})

	router.GET("/add/:name/:year", func(c *gin.Context) {
		name := c.Param("name")
		year, _ := strconv.Atoi(c.Param("year"))

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
			Name: name,
			Year: year,
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
			c.JSON(http.StatusOK, gin.H{
				"Added document": *item,
			})
		}

	})

	router.GET("/get/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		elem, found := db.Query(config.Namespace).
			Where("id", reindexer.EQ, id).
			Get()

		if found {
			item := elem.(*Item)
			c.JSON(http.StatusOK, gin.H{
				"Document": *item,
			})
		} else {
			c.String(http.StatusOK, "There is no such document.")
		}
	})

	router.GET("/getlist", func(c *gin.Context) {
		query := db.Query(config.Namespace).
			Sort("id", false).
			ReqTotal()

		iterator := query.Exec()
		defer iterator.Close()
		message := fmt.Sprintf("Found %v total documents. Documents:\n", iterator.TotalCount())
		c.String(http.StatusOK, message)

		// Iterate over results
		for iterator.Next() {
			elem := iterator.Object().(*Item)
			c.String(http.StatusOK, fmt.Sprintln(*elem))
		}

		if err := iterator.Error(); err != nil {
			panic(err)
		}
	})

	router.GET("/update/:id/:name/:year", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		name := c.Param("name")
		year, _ := strconv.Atoi(c.Param("year"))

		// Update document
		err := db.Upsert(config.Namespace, &Item{
			ID:   id,
			Name: name,
			Year: year,
		})
		if err != nil {
			panic(err)
		}

		// Request an updated document
		elem, found := db.Query(config.Namespace).
			Where("id", reindexer.EQ, id).
			Get()

		if found {
			item := elem.(*Item)
			c.JSON(http.StatusOK, gin.H{
				"Updated document": *item,
			})
		}
	})

	router.GET("/delete/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		err := db.Delete(config.Namespace, &Item{
			ID: id,
		})
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"delete task": id,
		})
	})

	router.Run(config.AppServerAddress)
}
