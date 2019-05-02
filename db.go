package main

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)


func ConnectDBClient()(*mongo.Client, error){

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()


	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		os.Getenv("DBadminID"),
		os.Getenv("DBadminPW"),
		os.Getenv("DBHOST"),
		os.Getenv("DBDatabase"),
	)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err !=nil{
		log.Error(fmt.Errorf("error occur in init of dataset handler\n"))
		log.Fatal(err.Error())
	}
	err = client.Connect(ctx)
	if err != nil {
		// FATAL Error : Fail to connect DB
		fmt.Printf("FATAL Error :  mongo client couldn't connect with background context %s", err.Error())
		log.Fatal(err.Error())
	}
	return client, err

}
