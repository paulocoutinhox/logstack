package datasource

import (
	"github.com/prsolucoes/logstack/models/domain"
	"gopkg.in/ini.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoDBDataSource struct {
	Session    *mgo.Session
	Host       string
	Database   string
	Collection string
}

type MongoDBLog struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Token     string        `bson:"token"`
	Type      string        `bson:"type"`
	Message   string        `bson:"message"`
	CreatedAt string        `bson:"created_at"`
}

func (This *MongoDBDataSource) Name() string {
	return "mongodb"
}

func (This *MongoDBDataSource) LoadConfig(config *ini.File) error {
	var err error

	serverSection, err := config.GetSection("server")

	if err != nil {
		This.Host = "localhost"
		This.Database = "logstack"
		This.Collection = "log"
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
	This.Database = dataBaseName
	This.Collection = "log"

	return nil
}

func (This *MongoDBDataSource) Connect() error {
	session, err := mgo.Dial("mongodb://" + This.Host)

	if err != nil {
		return err
	}

	This.Session = session
	return nil
}

func (This *MongoDBDataSource) Prepare() error {
	This.createCollections()
	This.createIndexes()
	return nil
}

func (This *MongoDBDataSource) InsertLog(newLog *models.Log) error {
	createAt := time.Now().UTC().Format("2006-01-02T15:04:05")
	newLog.CreatedAt = createAt

	mongoDBLog := &MongoDBLog{
		Token     : newLog.Token,
		Type      : newLog.Type,
		Message   : newLog.Message,
		CreatedAt : newLog.CreatedAt,
	}

	session := This.Session.Clone()
	defer session.Close()

	coll := session.DB(This.Database).C(This.Collection)
	err := coll.Insert(mongoDBLog)
	return err
}

func (This *MongoDBDataSource) LogList(token, message string, createdAt time.Time) ([]models.Log, error) {
	session := This.Session.Clone()
	defer session.Close()

	coll := session.DB(This.Database).C(This.Collection)

	var results []models.Log
	var conditions = bson.M{}

	conditions["token"] = token

	if message == "" {
		createdAtStr := createdAt.Format("2006-01-02T15:04:05")
		conditions["created_at"] = bson.M{"$gt": createdAtStr}
	} else {
		conditions["message"] = bson.RegEx{Pattern: message, Options: "i"}
	}

	var searchResult []MongoDBLog
	err := coll.Find(conditions).Sort("created_at").All(&searchResult)

	if err != nil {
		return nil, err
	}

	for _, mongoDBLog := range searchResult {
		log := models.Log{
			ID        : mongoDBLog.ID.Hex(),
			Token     : mongoDBLog.Token,
			Type      : mongoDBLog.Type,
			Message   : mongoDBLog.Message,
			CreatedAt : mongoDBLog.CreatedAt,
		}

		results = append(results, log)
	}

	return results, err
}

func (This *MongoDBDataSource) DeleteAllLogsByToken(token string) error {
	session := This.Session.Clone()
	defer session.Close()

	coll := session.DB(This.Database).C(This.Collection)

	_, err := coll.RemoveAll(bson.M{
		"token": token,
	})

	return err
}

func (This *MongoDBDataSource) LogStatsByType(token string) ([]models.LogStats, error) {
	session := This.Session.Clone()
	defer session.Close()

	coll := session.DB(This.Database).C(This.Collection)

	var results []models.LogStats

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

func (This *MongoDBDataSource) createCollections() {
	This.Session.DB(This.Database).C(This.Collection).Create(&mgo.CollectionInfo{DisableIdIndex: false, ForceIdIndex: true})
}

func (This *MongoDBDataSource) createIndexes() {
	This.Session.DB(This.Database).C(This.Collection).EnsureIndex(mgo.Index{Key: []string{"token"}, Unique: false, DropDups: true, Background: false, Sparse: true})
	This.Session.DB(This.Database).C(This.Collection).EnsureIndex(mgo.Index{Key: []string{"type"}, Unique: false, DropDups: true, Background: false, Sparse: true})
	This.Session.DB(This.Database).C(This.Collection).EnsureIndex(mgo.Index{Key: []string{"created_at"}, Unique: false, DropDups: true, Background: false, Sparse: true})
}
