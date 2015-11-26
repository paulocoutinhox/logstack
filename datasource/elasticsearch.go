package datasource

import (
	"github.com/prsolucoes/logstack/models/domain"
	"gopkg.in/ini.v1"
	"time"
	"gopkg.in/olivere/elastic.v3"
	"errors"
	"fmt"
	"log"
	"encoding/json"
)

type ElasticSearchDataSource struct {
	Client *elastic.Client
	Host   string
	Index  string
	Type   string
}

func (This *ElasticSearchDataSource) Name() string {
	return "elasticsearch"
}

func (This *ElasticSearchDataSource) LoadConfig(config *ini.File) error {
	var err error

	serverSection, err := config.GetSection("server")

	This.Type = "log"

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
	This.createMappings()
	return nil
}

func (This *ElasticSearchDataSource) InsertLog(newLog *models.Log) error {
	createAt := time.Now().UTC().Format("2006-01-02T15:04:05")
	newLog.CreatedAt = createAt

	_, err := This.Client.Index().
	Index(This.Index).
	Type(This.Type).
	BodyJson(newLog).
	Do()

	if err != nil {
		return err
	}

	return nil
}

func (This *ElasticSearchDataSource) LogList(token, message string, createdAt time.Time) ([]models.Log, error) {
	var results []models.Log

	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermsQuery("token", token))

	if message == "" {
		createdAtStr := createdAt.Format("2006-01-02T15:04:05")
		query = query.Must(elastic.NewRangeQuery("created_at").Gt(createdAtStr))
	} else {
		query = query.Must(elastic.NewMatchQuery("message", message))
	}

	source, err := query.Source()
	log.Printf("Query: %v", source)

	searchResult, err := This.Client.Search().Index(This.Index).Query(query).Do()

	if err != nil {
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != This.Index {
			return nil, errors.New(fmt.Sprintf("Expected SearchResult.Hits.Hit.Index = %q, but got %q", This.Index, hit.Index))
		}

		var log models.Log

		err := json.Unmarshal(*hit.Source, &log)

		if err != nil {
			return nil, err
		}

		log.ID = hit.Id
		results = append(results, log)
	}

	return results, err
}

func (This *ElasticSearchDataSource) DeleteAllLogsByToken(token string) error {
	return nil
}

func (This *ElasticSearchDataSource) LogStatsByType(token string) ([]models.LogStats, error) {
	query := elastic.NewMatchQuery("token", token)

	builder := This.Client.Search().Index(This.Index).Query(query)
	builder = builder.Aggregation("aggs-type", elastic.NewTermsAggregation().Field("type"))

	result, err := builder.Do()

	if err != nil {
		return nil, err
	}

	agg, found := result.Aggregations.Terms("aggs-type")

	if !found {
		return nil, errors.New("aggregation not found on result")
	}

	logStatList := []models.LogStats{}

	for _, bucket := range agg.Buckets {
		logStat := models.LogStats{}
		logStat.Type = bucket.Key.(string)
		logStat.Quantity = bucket.DocCount

		logStatList = append(logStatList, logStat)
	}

	return logStatList, nil
}

func (This *ElasticSearchDataSource) createIndex() {
	_, _ = This.Client.CreateIndex(This.Index).Do()
}

func (This *ElasticSearchDataSource) createMappings() {
	mapping := `
	{
        "log" : {
            "properties" : {
                "token" : {
                    "type" : "string",
                    "index" : "not_analyzed"
                },
                "type" : {
                    "type" : "string",
                    "index" : "not_analyzed"
                },
                "created_at" : {
                    "type" : "date",
                    "format" : "yyyy-MM-dd'T'HH:mm:ss",
                    "null_value": "now"
                }
            }
        }
    }
	`

	_, _ = This.Client.PutMapping().Index(This.Index).Type(This.Type).BodyString(mapping).Do()
}