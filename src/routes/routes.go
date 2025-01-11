package routes

import (
	"EeYenker/src/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "items.go.tmpl", gin.H{
			"cout":  5,
			"items": []string{"item1", "item2", "item3", "item4", "item5"},
		})
	})

	r.GET("/games", func(c *gin.Context) {
		c.HTML(http.StatusOK, "games.go.tmpl", gin.H{})
	})

}

func RegisterAPI(r *gin.Engine) {
	r.GET("/", controllers.Search)
}
