package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/prsolucoes/logstack/app"
	"github.com/prsolucoes/logstack/models/util"
	"log"
	"net/http"
)

type LogController struct{}

func (This *LogController) Register() {
	app.Server.Router.GET("/log/index", This.LogIndex)
	app.Server.Router.GET("/log/token", This.LogToken)
	app.Server.Router.GET("/log/statsByType", This.LogStatsByType)
	log.Println("LogController register : OK")
}

func (This *LogController) LogToken(c *gin.Context) {
	util.RenderTemplate(c.Writer, "log/token", nil)
}

func (This *LogController) LogIndex(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.Redirect(http.StatusMovedPermanently, "/log/token")
	}

	util.RenderTemplate(c.Writer, "log/index", map[string]string{"ContainerClass": "container-fluid", "WrapClass": "wrap-log", "FooterClass": "no-footer", "ShowLogMenu": "1"})
}

func (This *LogController) LogStatsByType(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.Redirect(http.StatusMovedPermanently, "/log/token")
	}

	util.RenderTemplate(c.Writer, "log/statsByType", map[string]string{"ContainerClass": "container-fluid", "WrapClass": "wrap-log", "FooterClass": "no-footer", "ShowLogMenu": "1"})
}
