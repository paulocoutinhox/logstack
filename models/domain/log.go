package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Log struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Token     string        `bson:"token" json:"token"`
	Type      string        `bson:"type" json:"type"`
	Message   string        `bson:"message" json:"message"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
