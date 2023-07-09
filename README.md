## Запуск образа doker

Пример

```bash
docker run -p9088:9088 -p6534:6534 -it reindexer/reindexer
```

## Использование API
Список HTTP-запросов
```
curl http://localhost:8080/
```

```
List of HTTP requests:
		/		Shows a list of commands or help for one command
		/getlist	Get a list of documents
		/get/:id	Get the document
		/add		Add a document
		/update		Update document
		/delete		Delete a document
```

### Пример создания документа

```
curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"John\", \"year\": 1998}" localhost:8080/add
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