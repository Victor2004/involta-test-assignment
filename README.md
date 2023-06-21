## Решение тестового задания involta

Программа представляет из себя консольное приложение которое подключается к серверу с Reindexer и позволяет создавать, редактировать, выводить информацию о списке имеющихся документов или заданного документа.

## Запуск образа doker

Пример

```bash
docker run -p9088:9088 -p6534:6534 -it reindexer/reindexer
```

## Использование приложения
Показать список команд или справку по одной команде
```bash
--help, -h
```

```
USAGE:
   reindex-db-loader [global options] command [command options] [arguments...]

COMMANDS:
   add, a       Add a document
   getlist, gl  Get a list of documents
   get, g       Get the document
   update, u    Update document
   delete, d    Delete a document
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### Пример создания документа

```bash
main add --name "Victor" --year 2004
```

или

```bash
main a -n Victor -y 2004
```

или

```bash
go run . add -n Victor --year 2004
```

## Конфигурационный файл

app.env находистся в папке configs

```env
SERVER_ADDRESS=cproto://localhost:6534
DATABASE_NAME=testdb
NAMESPACE=items
```
