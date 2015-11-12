package datasource

import (
	"github.com/prsolucoes/logstack/models/domain"
	"time"
	"gopkg.in/ini.v1"
)

type IDataSource interface {
	Name() string
	LoadConfig(config *ini.File) error
	Connect() error
	Prepare() error
	InsertLog(log *models.LogHistory) error
	LogList(token, message string, createdAt time.Time) ([]models.LogHistory, error)
	DeleteAllLogHistoryByToken(token string) error
	LogStatsByType(token string) ([]interface{}, error)
}
