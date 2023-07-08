## Запуск образа doker

Пример

```bash
docker run -p9088:9088 -p6534:6534 -it reindexer/reindexer
```

## Использование API
Список HTTP-запросов
```
http://localhost:8080/
```

```
List of HTTP requests:
		/		Shows a list of commands or help for one command
		/add/:name/:year		Add a document
		/get/:id		Get the document
		/getlist		Get a list of documents
		/update/:id/:name/:year		Update document
		/delete/:id		Delete a document
```

### Пример создания документа

```
http://localhost:8080/add/Victor/2004
```

```
{"Added document":{"ID":36,"Name":"Victor","Year":2004}}
```

## Конфигурационный файл

app.env находистся в папке configs

```env
REINDEXER_SERVER_ADDRESS=cproto://localhost:6534
DATABASE_NAME=testdb
NAMESPACE=items
APP_SERVER_ADDRESS=:8080
```