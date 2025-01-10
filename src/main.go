package main

import (
	"EeYenker/src/routes"

	"github.com/gin-gonic/gin" // Для маршрутизації
)

func main() {
	r := gin.Default()

	// Завантаження шаблонів

	r.LoadHTMLGlob("src/templates/**/*")

	r.Static("/static", "./src/static")
	r.Static("/assets", "./src/assets")

	// Реєстрація маршрутів
	routes.RegisterRoutes(r)
	routes.RegisterAPI(r)

	// Запуск сервера
	r.Run("0.0.0.0:8070")
}
