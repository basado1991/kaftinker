### kaftinker

Генерация ключей:
```sh
$ openssl genpkey -algorithm ed25519 -out private.pem
$ openssl pkey -in private.pem -pubout -out public.pem
```

Необходимые переменные окружения:
```sh
# Для блог-сервера
BS_ADDR=<адрес, куда биндиться серверу блога>
BS_ASSETS_PATH=<путь к папке assets/>
BS_TEMPLATES_PATH=<путь к папке templates/>
BS_PRIVATE_KEY=<путь к приватному PEM ed25519 ключу>
BS_PUBLIC_KEY=<путь к публичному PEM ed25519 ключу>
BS_SQLITE_DATABASE=<путь к базе данных SQL>
BS_PASSWORD_SALT=<соль для паролей>

KT_KAFKA_PARTITION=<партиция кафки для записи данных>
KT_KAFKA_ADDR=<адрес лидер-брокера>
KT_KAFKA_TOPIC=<топик>

# Для телеграм-консумера
TGC_TOKEN=<токен бота>
TGC_CHAT_ID=<айди чата, куда отправлять новые посты>
```
