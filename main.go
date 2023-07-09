package main

import (
	"fmt"
	"main/tools"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/restream/reindexer"
)

// Структура базы данных
type Item struct {
	ID   int64  `reindex:"id,,pk"`
	Name string `reindex:"name"`
	Year int    `reindex:"year,tree"`
}

func main() {
	// Инициализируем роутер Gin, используя Default.
	router := gin.Default()

	config, err := tools.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// Создаем подключение к серверу Reindexer
	db := reindexer.NewReindex(config.ReindexerServerAddress + "/" + config.DatabaseName)

	// Открываем нужное пространство имен
	db.OpenNamespace(config.Namespace, reindexer.DefaultNamespaceOptions(), Item{})

	// Используем функцию GET, чтобы связать метод GET HTTP и путь / с функцией обработчика.
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, `COMMANDS:
		/		Shows a list of commands or help for one command
		/getlist	Get a list of documents
		/get/:id	Get the document
		/add		Add a document
		/update		Update document
		/delete		Delete a document`)
	})

	// Используем функцию GET, чтобы связать метод GET HTTP и путь /getlist с функцией обработчика.
	router.GET("/getlist", func(c *gin.Context) {
		query := db.Query(config.Namespace).
			Sort("id", false).
			ReqTotal()

		iterator := query.Exec()
		defer iterator.Close()
		message := fmt.Sprintf("Found %v total documents. Documents:\n", iterator.TotalCount())
		c.String(http.StatusOK, message)

		// Итерация результатов
		for iterator.Next() {
			elem := iterator.Object().(*Item)
			c.String(http.StatusOK, fmt.Sprintln(*elem))
		}

		if err := iterator.Error(); err != nil {
			panic(err)
		}
	})

	// Используем функцию GET, чтобы связать метод GET HTTP и путь /get/X (вместо X нужно указать нужный id) с функцией обработчика.
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

	// Связываем метод POST по пути /add с функцией обработчика.
	router.POST("/add", func(c *gin.Context) {
		var request Item

		if err := c.BindJSON(&request); err != nil {
			return
		}

		// Находим идентификатор нового документа
		query := db.Query(config.Namespace).ReqTotal()
		iterator := query.Exec()
		defer iterator.Close()

		if err := iterator.Error(); err != nil {
			panic(err)
		}

		LastID := int64(iterator.TotalCount())

		// Добавляем документ в базу данных
		err := db.Upsert(config.Namespace, &Item{
			ID:   LastID,
			Name: request.Name,
			Year: request.Year,
		})

		if err != nil {
			panic(err)
		}

		// Запрос добавленного документа
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

	// Связываем метод PUT по пути /update с функцией обработчика.
	router.PUT("/update", func(c *gin.Context) {
		var request Item

		if err := c.BindJSON(&request); err != nil {
			return
		}

		// Обновляем документ
		err := db.Upsert(config.Namespace, &Item{
			ID:   request.ID,
			Name: request.Name,
			Year: request.Year,
		})
		if err != nil {
			panic(err)
		}

		// Запрос обновленного документа
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

	// Связываем метод DELETE по пути /delete с функцией обработчика.
	router.DELETE("/delete", func(c *gin.Context) {
		var request Item

		if err := c.BindJSON(&request); err != nil {
			return
		}

		// Удаляем документ
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

	// Запускаем веб сервис
	router.Run(config.AppServerPort)
}
