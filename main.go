package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"

	apidb "github.com/b0gochort/internal/api_db"
	"github.com/b0gochort/internal/handler"
	"github.com/b0gochort/internal/serivce"
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
	"github.com/valyala/fasthttp"
)

var db = reindexer.NewReindex("cproto://rx.web-gen.ru:6534/tinkoff")

func main() {

	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config; %s", err.Error())
	}

	type EmailConfig struct {
		From     string
		Password string
		SMTPHost string
		SMTPPort string
	}

	db := reindexer.NewReindex("cproto://rx.web-gen.ru:6534/tinkoff")
	if db.Status().Err != nil {
		panic(db.Status().Err)
	}

	apiDb := apidb.NewAPIDB(db)

	service := serivce.NewService(apiDb)

	handler := handler.NewHandler(service)

	server := fasthttp.Server{
		Handler: handler.InitRoutes, // Обработчик запросов
	}

	// Слушаем порт 8080 и обрабатываем запросы

	fmt.Println("ok")
	err := server.ListenAndServe(":8080")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}

}
