package app

import (
	"github.com/gin-gonic/gin"
	"github.com/prsolucoes/logstack/datasource"
	"log"
	"gopkg.in/ini.v1"
)

type WebServer struct {
	DS     datasource.IDataSource
	Router *gin.Engine
	Config *ini.File
	Port   int
}

var (
	Server *WebServer
)

func NewWebServer() *WebServer {
	server := new(WebServer)

	gin.SetMode(gin.ReleaseMode)
	server.Router = gin.New()
	server.Router.Use(gin.Recovery())

	return server
}

func (This *WebServer) CreateBasicRoutes() {
	This.Router.Static("/static", "resources/static")
	log.Println("Router creation : OK")
}

func (This *WebServer) LoadConfiguration() {
	config, err := ini.Load([]byte(""), "config.ini")

	if err == nil {
		This.Config = config

		serverSection, err := config.GetSection("server")

		if err != nil {
			This.Port = 8080
		}

		port, err := serverSection.Key("port").Int()

		if err != nil || port <= 0 {
			port = 8080
		}

		This.Port = port

		log.Println("Configuration file load : OK")
	} else {
		log.Fatalf("Configuration file load error : %s", err.Error())
	}
}

func (This *WebServer) CreateDataSource() {
	var err error

	serverSection, err := This.Config.GetSection("server")
	dsnameKey := serverSection.Key("dsname")
	dsname := dsnameKey.String()

	// create datasource from config
	if dsname == "mongodb" {
		This.DS = new(datasource.MongoDataSource)
	} else {
		log.Fatal("No datasource defined or is invalid")
	}

	log.Println("DataSource defined : " + This.DS.Name())

	// load datasource config
	err = This.DS.LoadConfig(This.Config)

	if err == nil {
		log.Println("DataSource config load : OK")
	} else {
		log.Fatalf("DataSource config load error : %s", err)
	}

	// connect to datasource server
	err = This.DS.Connect()

	if err == nil {
		log.Println("DataSource connect : OK")
	} else {
		log.Fatalf("DataSource connect error : %s", err)
	}

	// prepare datasource
	err = This.DS.Prepare()

	if err == nil {
		log.Println("DataSource prepare : OK")
	} else {
		log.Fatalf("DataSource prepare error : %s", err)
	}

	// finish datasource creation
	log.Println("DataSource creation : OK")
}

func (This *WebServer) Start() {
	log.Printf("Server started on port %v : OK", This.Port)
	This.Router.Run(":8080")
}
