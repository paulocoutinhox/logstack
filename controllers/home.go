package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/prsolucoes/logstack/app"
	"github.com/prsolucoes/logstack/models/util"
	"log"
)

type HomeController struct{}

func (This *HomeController) Register() {
	app.Server.Router.GET("/", This.HomeIndex)
	log.Println("HomeController register : OK")
}

func (This *HomeController) HomeIndex(c *gin.Context) {
	util.RenderTemplate(c.Writer, "home/index", nil)
}
