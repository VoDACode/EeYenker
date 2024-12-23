package main

import (
	"github.com/gin-gonic/gin"           // Для маршрутизації
	"EeYenker/src/routes"
)

func main() {
	r := gin.Default()

	// Завантаження шаблонів
	
	r.LoadHTMLGlob("src/templates/**/*")

	r.Static("/static", "./src/static")

	// Реєстрація маршрутів
	routes.RegisterRoutes(r)
	routes.RegisterAPI(r)

	// Запуск сервера
	r.Run("0.0.0.0:8080")
}
