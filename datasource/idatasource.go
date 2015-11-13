package datasource

import (
	"github.com/prsolucoes/logstack/models/domain"
	"gopkg.in/ini.v1"
	"time"
)

type IDataSource interface {
	Name() string
	LoadConfig(config *ini.File) error
	Connect() error
	Prepare() error
	InsertLog(log *models.Log) error
	LogList(token, message string, createdAt time.Time) ([]models.Log, error)
	DeleteAllLogsByToken(token string) error
	LogStatsByType(token string) ([]models.LogStats, error)
}
