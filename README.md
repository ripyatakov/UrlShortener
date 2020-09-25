# UrlShortener
Программа-сервис для сокращения ссылок. На вход принимаются URL с JSON файлом для получения сжатой ссылки или кастомной сжатой ссылки, либо только URL для переадресации 

## Настройка БД
Чтобы создать нужную таблицу базы данных необходимо выполнить следующий код в SQL Editor:
```sql
CREATE TABLE `urls` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `surl` char(101) DEFAULT NULL,
  `lurl` varchar(101) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;'
```
### Сокращение ссылок
Для сокращения ссылки пользователь должен передать отправить запрос следующего вида
`http://localhost:8000/url` параметром `POST` и JSON вида `{"lurl": "https://www.avito.ru"}`
При сокращении используется алгоритм base64 из стандартной библиотеки `encoding` с некоторой модификацией
### Получение кастомной ссылки
Для получения ссылки пользователь должен передать отправить запрос следующего вида
`http://localhost:8000/сurl` параметром `POST` и JSON вида `{"lurl": "https://www.avito.ru", "surl":"avito"}`
При сокращении используется алгоритм base64 из стандартной библиотеки `encoding` с некоторой модификацией
### Переадресация
Для переадресации пользователь должен перейти по ссылки вида
`http://localhost:8000/avito`