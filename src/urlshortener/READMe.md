### Urlsshortener

Сервис для обработки коротких ссылок

#### Команда для запуска

``
urlsshortener -f config.json
``

#### Пример файла конфигурации

```json
{
    "port": 3000,
    "default_path": "https://golang.org/",
    "paths": {
        "/go-youtube": "https://youtube.com/",
        "/go-gh": "https://github.com."
    }
}
```

#### Пример логов работы сервиса
```
2021/02/22 17:23:59 Starting server on port: 3000
2021/02/22 17:24:07 Moved from / to https://golang.org/
2021/02/22 17:25:21 Moved from /go-youtube to https://youtube.com/
```
