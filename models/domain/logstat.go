package models

type LogStats struct {
	Type     string   `bson:"type" json:"type"`
	Quantity int64    `bson:"quantity" json:"quantity"`
}
