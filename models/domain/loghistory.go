package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type LogHistory struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Token     string        `bson:"token" json:"token"`
	Type      string        `bson:"type" json:"type"`
	Message   string        `bson:"message" json:"message"`
	CreatedAt time.Time     `bson:"createdAt" json:"created_at"`
}
