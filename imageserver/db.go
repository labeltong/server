package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type key string

const (
	hostKey     = key("DBHOST")
	usernameKey = key("DBadminID")
	passwordKey = key("DBadminPW")
	databaseKey = key("DBDatabase")
)

type MongoDBInstance struct {
	DBClient *mongo.Client
}
type MongoDBCollection struct {
	CollectionName string
	Collection     *mongo.Collection
}

func InitMongoDBClient() (mdbi MongoDBInstance, err error) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, hostKey, os.Getenv("DBHOST"))
	ctx = context.WithValue(ctx, usernameKey, os.Getenv("DBadminID"))
	ctx = context.WithValue(ctx, passwordKey, os.Getenv("DBadminPW"))
	ctx = context.WithValue(ctx, databaseKey, os.Getenv("DBDatabase"))

	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		ctx.Value(usernameKey).(string),
		ctx.Value(passwordKey).(string),
		ctx.Value(hostKey).(string),
		ctx.Value(databaseKey).(string),
	)

	mdbi.DBClient, err = mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		// FATAL Error : Fail to connect DB
		fmt.Printf("FATAL Error : Fail to connect DB with %s", err.Error())
		log.Fatal(err.Error())
	}

	err = mdbi.DBClient.Connect(ctx)
	if err != nil {
		// FATAL Error : Fail to connect DB
		fmt.Printf("FATAL Error :  mongo client couldn't connect with background context %s", err.Error())
		log.Fatal(err.Error())
	}

	return mdbi, err
}

func InitDatabase(collName string, dbname string) (col MongoDBCollection, err error) {
	mdbi, err := InitMongoDBClient()
	if err != nil {
		// FATAL Error : Fail to connect DB
		fmt.Printf("FATAL Error : Fail to connect DB with %s", err.Error())
		log.Fatal(err.Error())
	}
	col.CollectionName = collName
	col.Collection = mdbi.DBClient.Database(dbname).Collection(collName)
	return col, err

}

func (mdbc *MongoDBCollection) InsertData(data datatolabel) (err error) {
	insertResult, err := mdbc.Collection.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Println("FATAL Error: Fail to insert Data")
		log.Fatal(err)
		return err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil

}

func (mdbc *MongoDBCollection) RemoveData(dataId string) (err error) {
	deleteResult, err := mdbc.Collection.DeleteOne(context.TODO(), bson.M{"fileid": dataId})
	if err != nil {
		fmt.Println("FATAL Error: Fail to delete Data")
		log.Fatal(err)
		return err
	}
	fmt.Println("Deleted a single document: ", deleteResult.DeletedCount)
	return nil
}

func (mdbc *MongoDBCollection) ReadData(dataId string) (data *datatolabel, err error) {
	err = mdbc.Collection.FindOne(context.TODO(), bson.M{"fileid": dataId}).Decode(&data)
	if err != nil {
		fmt.Println("FATAL Error: Fail to read Data")
		log.Fatal(err)

	}
	fmt.Println("Read a single document: ")
	return data, nil
}
