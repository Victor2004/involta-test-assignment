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
		/getlist	Get a list of documents
		/get/:id	Get the document
		/add		Add a document
		/update		Update document
		/delete		Delete a document`)
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

	// curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"John\", \"year\": 1998}" localhost:8080/add
	router.POST("/add", func(c *gin.Context) {
		var request Item

		if err := c.BindJSON(&request); err != nil {
			return
		}

		// c.JSON(http.StatusCreated, gin.H{name.Name: name.Year})
		// name := c.Param("name")
		// year, _ := strconv.Atoi(c.Param("year"))

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
			Name: request.Name,
			Year: request.Year,
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

	// curl -X PUT -H "Content-Type: application/json" -d "{\"id\": 20, \"name\":\"John\", \"year\": 1998}" localhost:8080/update
	router.PUT("/update", func(c *gin.Context) {
		var request Item

		if err := c.BindJSON(&request); err != nil {
			return
		}

		// Update document
		err := db.Upsert(config.Namespace, &Item{
			ID:   request.ID,
			Name: request.Name,
			Year: request.Year,
		})
		if err != nil {
			panic(err)
		}

		// Request an updated document
		elem, found := db.Query(config.Namespace).
			Where("id", reindexer.EQ, request.ID).
			Get()

		if found {
			item := elem.(*Item)
			c.JSON(http.StatusOK, gin.H{
				"Updated document": *item,
			})
		}
	})

	// curl -X DELETE -H "Content-Type: application/json" -d "{\"id\": 43}" localhost:8080/delete
	router.DELETE("/delete", func(c *gin.Context) {
		var request Item

		if err := c.BindJSON(&request); err != nil {
			return
		}

		err := db.Delete(config.Namespace, &Item{
			ID: request.ID,
		})
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"delete task": request.ID,
		})
	})

	router.Run(config.AppServerAddress)
}
