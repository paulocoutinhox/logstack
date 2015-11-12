package datasource

import (
	"github.com/prsolucoes/logstack/models/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/ini.v1"
)

type MongoDataSource struct {
	Session      *mgo.Session
	Host         string
	DataBaseName string
}

func (This *MongoDataSource) Name() string {
	return "mongods"
}

func (This *MongoDataSource) LoadConfig(config *ini.File) error {
	var err error

	serverSection, err := config.GetSection("server")

	if err != nil {
		This.Host = "localhost"
		This.DataBaseName = "logstack"
		return nil
	}

	hostKey := serverSection.Key("host")
	dataBaseNameKey := serverSection.Key("databasename")

	host := hostKey.String()
	dataBaseName := dataBaseNameKey.String()

	if host == "" {
		host = "localhost"
	}

	if dataBaseName == "" {
		host = "logstack"
	}

	This.Host = host
	This.DataBaseName = dataBaseName

	return nil
}

func (This *MongoDataSource) Connect() error {
	session, err := mgo.Dial("mongodb://" + This.Host)

	if err != nil {
		return err
	}

	This.Session = session
	return nil
}

func (This *MongoDataSource) Prepare() error {
	This.createCollections()
	This.createIndexes()
	return nil
}

func (This *MongoDataSource) InsertLog(log *models.LogHistory) error {
	session := This.Session.Clone()
	defer session.Close()

	coll := session.DB(This.DataBaseName).C("loghistory")
	err := coll.Insert(log)
	return err
}

func (This *MongoDataSource) LogList(token, message string, createdAt time.Time) ([]models.LogHistory, error) {
	session := This.Session.Clone()
	defer session.Close()

	coll := session.DB(This.DataBaseName).C("loghistory")

	var results []models.LogHistory
	var conditions = bson.M{}

	conditions["token"] = token

	if message == "" {
		conditions["created_at"] = bson.M{"$gt": createdAt}
	} else {
		conditions["message"] = bson.RegEx{Pattern: message, Options: "i"}
	}

	err := coll.Find(conditions).Sort("createdAt").All(&results)
	return results, err
}

func (This *MongoDataSource) DeleteAllLogHistoryByToken(token string) error {
	session := This.Session.Clone()
	defer session.Close()

	coll := session.DB(This.DataBaseName).C("loghistory")

	_, err := coll.RemoveAll(bson.M{
		"token": token,
	})

	return err
}

func (This *MongoDataSource) LogStatsByType(token string) ([]interface{}, error) {
	session := This.Session.Clone()
	defer session.Close()

	coll := session.DB(This.DataBaseName).C("loghistory")

	var results []interface{}

	pipe := coll.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"token": token,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$type",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"type": bson.M{
					"$toLower": "$_id",
				},
				"quantity": "$count",
				"_id":      0,
			},
		},
	})

	err := pipe.All(&results)
	return results, err
}

func (This *MongoDataSource) createCollections() {
	This.Session.DB(This.DataBaseName).C("loghistory").Create(&mgo.CollectionInfo{DisableIdIndex: false, ForceIdIndex: true})
}

func (This *MongoDataSource) createIndexes() {
	This.Session.DB(This.DataBaseName).C("loghistory").EnsureIndex(mgo.Index{Key: []string{"token"}, Unique: false, DropDups: true, Background: false, Sparse: true})
	This.Session.DB(This.DataBaseName).C("loghistory").EnsureIndex(mgo.Index{Key: []string{"type"}, Unique: false, DropDups: true, Background: false, Sparse: true})
	This.Session.DB(This.DataBaseName).C("loghistory").EnsureIndex(mgo.Index{Key: []string{"created_at"}, Unique: false, DropDups: true, Background: false, Sparse: true})
}
