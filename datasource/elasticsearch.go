package datasource

import (
	"github.com/prsolucoes/logstack/models/domain"
	"gopkg.in/ini.v1"
	"time"
	"gopkg.in/olivere/elastic.v3"
)

type ElasticSearchDataSource struct {
	Client *elastic.Client
	Host   string
	Index  string
}

func (This *ElasticSearchDataSource) Name() string {
	return "elasticsearch"
}

func (This *ElasticSearchDataSource) LoadConfig(config *ini.File) error {
	var err error

	serverSection, err := config.GetSection("server")

	if err != nil {
		This.Host = "localhost"
		This.Index = "logstack"
		return nil
	}

	hostKey := serverSection.Key("dshost")
	dataBaseNameKey := serverSection.Key("dscontainer")

	host := hostKey.String()
	dataBaseName := dataBaseNameKey.String()

	if host == "" {
		host = "localhost"
	}

	if dataBaseName == "" {
		host = "logstack"
	}

	This.Host = host
	This.Index = dataBaseName

	return nil
}

func (This *ElasticSearchDataSource) Connect() error {
	client, err := elastic.NewClient(elastic.SetURL(This.Host))

	if err != nil {
		return err
	}

	This.Client = client
	return nil
}

func (This *ElasticSearchDataSource) Prepare() error {
	This.createIndex()
	return nil
}

func (This *ElasticSearchDataSource) InsertLog(log *models.LogHistory) error {
	_, err := This.Client.Index().
	Index(This.Index).
	Type("log").
	BodyJson(log).
	Do()

	if err != nil {
		return err
	}

	return nil
}

func (This *ElasticSearchDataSource) LogList(token, message string, createdAt time.Time) ([]models.LogHistory, error) {
	return nil, nil
}

func (This *ElasticSearchDataSource) DeleteAllLogHistoryByToken(token string) error {
	return nil
}

func (This *ElasticSearchDataSource) LogStatsByType(token string) ([]interface{}, error) {
	return nil, nil
}

func (This *ElasticSearchDataSource) createIndex() {
	_, _ = This.Client.CreateIndex(This.Index).Do()
}
