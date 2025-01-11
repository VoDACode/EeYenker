package main

import (
	"EeYenker/src/routes"
	"fmt"

	"github.com/fatih/color" // Для кольорового тексту
	"github.com/gin-gonic/gin"
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

	// Виведення адреси сервера в консоль із підсвіткою
	port := "8070"
	url := fmt.Sprintf("http://localhost:%s", port)

	// Використовуємо кольоровий текст
	c := color.New(color.FgHiGreen, color.Bold) // Високий зелений текст, жирний
	c.Printf("Сервер запущено на: %s\n", url)

	// Запуск сервера
	r.Run("0.0.0.0:" + port)
}
