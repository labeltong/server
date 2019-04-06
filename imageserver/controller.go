package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

/*
/labelme/{id} [GET]    returns a specific person from the database
/labelme      [POST]   creates a new person in the database
/labelme      [DELETE] deletes a person from the database
/labelme/{id} [PUT]    updates the document of a person
*/

//CREATE

func (d *dataset) InitDataset(setname string, isbait bool, numofdata int, answertype string) (err error){

	d.datacoll ,err = InitDatabase(setname, "labeltong")
	d.DatasetName = setname
	d.isBait= isbait
	d.numofdata = numofdata
	d.answerType = answertype

	if err != nil{
		fmt.Println("FATAL Error: Fail to init database")
	}
	return err
}

func (d *dataset) InsertData(data datatolabel) (err error) {

	insertResult, err := d.datacoll.Collection.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Println("FATAL Error: Fail to insert Data")
		log.Fatal(err)
		return err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil

}

//READ

func (d *dataset) ReadData(dataId string) (data datatolabel, err error) {
	err = d.datacoll.Collection.FindOne(context.TODO(), bson.M{"fileid": dataId}).Decode(&data)
	if err != nil {
		fmt.Println("FATAL Error: Fail to read Data")
		log.Fatal(err)

	}
	fmt.Println("Read a single document: ")
	return data, nil
}

func (d *dataset) ReadAllData() ([]datatolabel, error) {
	var results []datatolabel

	cur, err := d.datacoll.Collection.Find(context.TODO(), nil, options.Find())
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem datatolabel
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	if err != nil {
		fmt.Println("FATAL Error: Fail to read Data")
		log.Fatal(err)

	}
	fmt.Println("Read a single document: ")
	return results, nil

}

//UPDATE
func (d *dataset) UpdateData(dataId string, updated datatolabel) (data datatolabel, err error) {
	retval := d.datacoll.Collection.FindOneAndReplace(context.TODO(), bson.M{"fileid": dataId}, updated).Decode(&data)

	fmt.Println("Replace origin : ")
	fmt.Println(retval)
	fmt.Println("Read a single document: ")
	return data, nil
}

//DELETE
func (d *dataset) DeleteData(dataId string) (err error) {
	deleteResult, err := d.datacoll.Collection.DeleteOne(context.TODO(), bson.M{"fileid": dataId})
	if err != nil {
		fmt.Println("FATAL Error: Fail to delete Data")
		log.Fatal(err)
		return err
	}
	fmt.Println("Deleted a single document: ", deleteResult.DeletedCount)
	return nil
}
