package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func writeInternalError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
}

func indexView(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func saveMessageView(c *gin.Context) {
	message := c.PostForm("message")
	key := keyBuilder.Get()
	err := keeper.Set(key, message)
	if err != nil {
		writeInternalError(c)
		return
	}
	c.HTML(http.StatusOK, "key.html", gin.H{"key": fmt.Sprintf("http://%s/%s", c.Request.Host, key)})
}

func readMessageView(c *gin.Context) {
	key := c.Param("key")
	msg, err := keeper.Get(key)
	if err != nil {
		if err.Error() == NotFoundError {
			c.HTML(http.StatusNotFound, "404.html", gin.H{})
			return
		}
		writeInternalError(c)
		return
	}

	err = keeper.Clean(key)
	if err != nil {
		writeInternalError(c)
		return
	}

	c.HTML(http.StatusOK, "message.html", gin.H{"message": msg})

}

func getRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles(
		"templates/index.html",
		"templates/key.html",
		"templates/message.html",
		"templates/404.html",
		"templates/500.html",
	)
	router.GET("/", indexView)
	router.POST("/", saveMessageView)
	router.GET("/:key", readMessageView)
	return router
}

func main() {
	router := getRouter()
	err := router.Run("localhost:9000")
	if err != nil {
		return
	}
}
