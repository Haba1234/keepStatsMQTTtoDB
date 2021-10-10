package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Haba1234/keepStatsMQTTtoDB/internal/app"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// DB структура сервера.
type DB struct {
	bucket   string
	org      string
	token    string
	url      string
	client   influxdb2.Client
	pointsCh <-chan app.Point
}

const timeout = 10 // Пауза между проверками готовности БД к работе (при первом запуске БД запускается не мгновенно).

// NewStorage конструктор.
func NewStorage(cfg app.StorageConf) *DB {
	return &DB{
		bucket: cfg.Bucket,
		org:    cfg.Org,
		token:  cfg.Token,
		url:    cfg.URL,
	}
}

// TODO подумать над cancel. Вызов при N кол-ве безуспешных попыток подключиться к БД.

// Start функция запускает сервис.
func (db *DB) Start(ctx context.Context, pointsCh <-chan app.Point) error {
	db.pointsCh = pointsCh

test:
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(timeout * time.Second):
			db.client = influxdb2.NewClient(db.url, db.token)
			ok, err := db.client.Ready(ctx)
			if err != nil {
				log.Println("DB not ready:", err)
				break
			}
			if ok {
				log.Println("DB ready!")
				break test
			}
		}
	}

	writeAPI := db.client.WriteAPI(db.org, db.bucket)

	errorsCh := writeAPI.Errors() // TODO проверить необходимость закрытия канала. возможно гонка.
	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			fmt.Printf("write error: %s\n", err.Error())
		}
	}()

	go db.writePoint(ctx, writeAPI)

	return nil
}

// Stop функция закрывает соединение с БД.
func (db *DB) Stop() error {
	db.client.Close()
	return nil
}

func (db *DB) writePoint(ctx context.Context, writeAPI api.WriteAPI) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case data := <-db.pointsCh:
			writeAPI.WritePoint(preparePoint(data))
		}
	}
	writeAPI.Flush()
}

func preparePoint(data app.Point) *write.Point {
	log.Println("Перед записью data:", data)
	return influxdb2.NewPoint(data.Measurement,
		data.Tag,
		map[string]interface{}{
			data.Field: data.Value,
		},
		time.Now())
}
