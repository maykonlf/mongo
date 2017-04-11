package mongo

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

func NewClientStore(config *Config) (cs *MongoClientStore) {
	cs = &MongoClientStore {config: config, collectionName: "clients"}

	session, err := mgo.Dial(cs.config.URL)
	if err != nil {
		return
	}
	cs.session = session

	err = cs.c(cs.collectionName).EnsureIndex(mgo.Index{
		Key: []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	})
	if err != nil {
		return
	}

	return cs
}

// MongoClientStore client information store
type MongoClientStore struct {
	config *Config
	collectionName string
	session *mgo.Session
}

// GetByID according to the ID for the client information
func (cs *MongoClientStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	session := cs.session.Copy()

	var client models.Client
	err = session.DB(cs.config.DB).C(cs.collectionName).Find(bson.M{"id": id}).One(&client)
	cli = &client
	if err != nil {
		err = errors.New("not found")
	}

	return
}

// Set set client information
func (cs *MongoClientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	session := cs.session.Copy()
	_, err = session.DB(cs.config.DB).C(cs.collectionName).Upsert(bson.M{
		"id": id,
	},
	&models.Client{
		ID: cli.GetID(),
		Secret: cli.GetSecret(),
		Domain: cli.GetDomain(),
		UserID: cli.GetUserID(),
	})

	return
}

func (cs *MongoClientStore) c(name string) *mgo.Collection {
	return cs.session.DB(cs.config.DB).C(name)
}
