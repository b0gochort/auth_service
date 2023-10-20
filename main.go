package main

import (
	"fmt"

	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
	"github.com/valyala/fasthttp"
)

var db = reindexer.NewReindex("cproto://rx.web-gen.ru:6534/tinkoff")

func handleRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/hello":
		handleHello(ctx)
	default:
		handleNotFound(ctx)
	}
	// Отправляем ответ клиенту
}

// Обработка корневого маршрута "/"
func handleRoot(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Добро пожаловать на главную страницу!")
}

// Обработка маршрута "/hello"
func handleHello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Привет, мир!")
}

// Обработка неизвестных маршрутов
func handleNotFound(ctx *fasthttp.RequestCtx) {
	ctx.Error("Страница не найдена", fasthttp.StatusNotFound)
}
func main() {
	if db.Status().Err != nil {
		panic(db.Status().Err)
	}

	server := fasthttp.Server{
		Handler: handleRequest, // Обработчик запросов
	}

	// Слушаем порт 8080 и обрабатываем запросы

	fmt.Println("ok")
	err := server.ListenAndServe(":8080")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}

}
