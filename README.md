MQTT клиент

## Описание
Может подписываться на несколько MQTT брокеров и сохраняет полученные данные в базу данных InfluxDB

### Конфигурация
Через командную строку задается путь к конфигурационному файлу

`service-name -config ./conf.toml`

### Сборка и запуск
```
git clone https://github.com/Haba1234/keepStatsMQTTtoDB.git
cd keepStatsMQTTtoDB
make build
make run
```
### Docker
```
sudo docker build -t service .
sudo docker run -d --name service -p 1886:1886 service
```
### Docker Compose
Разворачивает на хосте контейнеры с MQTT клиентом и базой InfluxDB.
Автоматически выполняется первоначальная настройка БД. 
Дополнительные данные из файла .env
```
sudo docker-compose up -d --build
```
