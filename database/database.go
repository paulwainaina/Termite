package Db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
	Ctx    context.Context
	Err    error
}

func (db *Database) Connect() {
	db.Client, db.Err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB")))
	if db.Err != nil {
		log.Fatal(db.Err)
	}
	db.Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	db.Err = db.Client.Connect(db.Ctx)
	if db.Err != nil {
		log.Fatal(db.Err)
	}
	db.Err = db.Client.Ping(db.Ctx, nil)
	if db.Err != nil {
		log.Fatal(db.Err)
	}
}

func (db *Database) DisConnect() {
	db.Client.Disconnect(db.Ctx)
}

func (db *Database) InsertOne(database string, col string, doc interface{}) (interface{}, error) {
	_col := db.Client.Database(database).Collection(col)
	var result interface{}
	db.Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	result, db.Err = _col.InsertOne(db.Ctx, doc)
	return result, db.Err
}

func (db *Database) QueryMany(database string, col string, filter interface{}) ([]bson.D, error) {
	_col := db.Client.Database(database).Collection(col)
	db.Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	cusor, err := _col.Find(db.Ctx, filter)
	result := []bson.D{}
	if err == nil {
		err = cusor.All(db.Ctx, &result)
	}
	db.Err = err
	return result, db.Err
}
