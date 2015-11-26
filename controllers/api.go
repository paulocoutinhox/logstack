package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/prsolucoes/gowebresponse"
	"github.com/prsolucoes/logstack/app"
	"github.com/prsolucoes/logstack/models/domain"
	"log"
	"strings"
	"time"
)

type APIController struct{}

func (This *APIController) Register() {
	app.Server.Router.POST("/api/log/add", This.APILogAdd)
	app.Server.Router.GET("/api/log/list", This.APILogList)
	app.Server.Router.GET("/api/log/deleteAll", This.APILogDeleteAll)
	app.Server.Router.GET("/api/log/statsByType", This.APILogStatsByType)
	log.Println("APIController register : OK")
}

func (This *APIController) APILogAdd(c *gin.Context) {
	log := &models.Log{}
	log.Token = c.PostForm("token")
	log.Type = strings.ToLower(c.PostForm("type"))
	log.Message = c.PostForm("message")
	log.CreatedAt = ""

	err := app.Server.DS.InsertLog(log)

	response := new(gowebresponse.WebResponse)

	if err == nil {
		response.Success = true
		response.Message = ""
	} else {
		response.Success = false
		response.Message = "log-add-error"
		response.AddDataError("error", err.Error())
	}

	c.JSON(200, response)
}

func (This *APIController) APILogList(c *gin.Context) {
	token := c.Query("token")
	createdAt, _ := time.Parse("2006-01-02T15:04:05.999", c.Query("created_at"))
	message := c.Query("message")

	list, err := app.Server.DS.LogList(token, message, createdAt)

	response := new(gowebresponse.WebResponse)

	if err == nil {
		response.Success = true
		response.Message = ""
		response.AddData("list", list)
	} else {
		response.Success = false
		response.Message = "log-list-error"
		response.AddDataError("error", err.Error())
	}

	c.JSON(200, response)
}

func (This *APIController) APILogDeleteAll(c *gin.Context) {
	token := c.Query("token")

	err := app.Server.DS.DeleteAllLogsByToken(token)

	response := new(gowebresponse.WebResponse)

	if err == nil {
		response.Success = true
		response.Message = ""
	} else {
		response.Success = false
		response.Message = "log-delete-all-error"
		response.AddDataError("error", err.Error())
	}

	c.JSON(200, response)
}

func (This *APIController) APILogStatsByType(c *gin.Context) {
	token := c.Query("token")

	list, err := app.Server.DS.LogStatsByType(token)

	response := new(gowebresponse.WebResponse)

	if err == nil {
		response.Success = true
		response.Message = ""
		response.AddData("list", list)
	} else {
		response.Success = false
		response.Message = "log-stats-by-type-error"
		response.AddDataError("error", err.Error())
	}

	c.JSON(200, response)
}
