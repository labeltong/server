package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)
/*
/dataset/list [GET] Login function, Oauth from client required
/dataset/list/{dsid}/get [GET] Logout function, Oauth from client required
/dataset/list/{dsid}/{id}/ans [POST] Check if user is authenticated not yet implemented
/dataset/list/{dsid}/{id}/info [GET] Check if user is authenticated not yet implemented

*/

var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`



var client *mongo.Client

func DatasetInitSubrouter(r *mux.Router)  {
	ret := r.PathPrefix("/dataset").Subrouter()

	ret.HandleFunc("/list", getAlldatasets).Methods("GET")
	ret.HandleFunc("/list/{dsid}/get", getDataFromDataset).Methods("GET")
	//ret.HandleFunc("/list/{dsid}/{id}/ans",postAnswerToDataset).Methods("POST")
	//ret.HandleFunc("/list/{dsid}/info",getDataSetInfo).Methods("GET")


}

func init(){
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()



	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		os.Getenv("DBadminID"),
		os.Getenv("DBadminPW"),
		os.Getenv("DBHOST"),
		os.Getenv("DBDatabase"),
	)


	client, err = mongo.NewClient(options.Client().ApplyURI(uri))
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


}


func getAlldatasets(w http.ResponseWriter, r *http.Request){
	r.Header.Set("Content-Type", "application/json")


	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")



	var results []dataset
	collection := client.Database(os.Getenv("DBDatabase")).Collection("dataset_list")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		if err !=nil{
			log.Error(fmt.Errorf("error occur Writing internal server error\n"))
			log.Fatal(err.Error())
		}
		return
	}

	defer func(){
		err2 := cursor.Close(ctx)
		if err2 != nil{
			log.Error(fmt.Errorf("error occur closing db cursor\n"))
			log.Fatal(err.Error())
		}
	}()
	for cursor.Next(ctx) {
		var result dataset
		err = cursor.Decode(&result)
		if err != nil{
			log.Error(fmt.Errorf("error occur in Decode data to dataset struct\n"))
			log.Fatal(err.Error())
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		if err !=nil{
			log.Error(fmt.Errorf("error occur Writing internal server error\n"))
			log.Fatal(err.Error())
		}
		return
	}
	err = json.NewEncoder(w).Encode(results)
	if err != nil{
		log.Error(fmt.Errorf("error occur in endocing json results\n"))
		log.Fatal(err.Error())
	}
	fmt.Println(r.Host + "Request dataset list "+ time.Now().UTC().String())

}


func getDataFromDataset(w http.ResponseWriter, r *http.Request){


	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	dsName := params["dsid"]
	var dsresult dataset
	var result datatolabel
	templateman := false
	collection := client.Database(os.Getenv("DBDatabase")).Collection("dataset_list")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	res := collection.FindOne(ctx, bson.M{"datasetdbname":dsName})
	err := res.Decode(&dsresult)
	if err != nil{
		log.Error(fmt.Errorf("error occur Decoding found dataset struct\n"))
		log.Fatal(err.Error())
	}
	numofdata := dsresult.Numofdata
	rand.Seed(time.Now().UnixNano())
	r1 := rand.Intn(numofdata - 1) + 1





	collection = client.Database(os.Getenv("DBDatabase")).Collection(dsresult.Datasetdbname)
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	res = collection.FindOne(ctx, bson.D{{"file_num" , r1}})

	err = res.Decode(&result)

	if err != nil{
		log.Error(fmt.Errorf("error occur Decoding found datatolabel struct\n"))
		log.Fatal(err.Error())
	}
	b, err := os.Open(result.Path)
	if err != nil{
		log.Error(fmt.Errorf("error occur loading images\n"))
		log.Fatal(err.Error())
	}
	defer func(){
		e := b.Close()
		if e!=nil{
			log.Error(fmt.Errorf("error occur closing files\n"))
			log.Fatal(err.Error())
		}
	}()// just pass the file name
	fInfo, _ := b.Stat()
	buf := make([]byte, fInfo.Size())

	fReader := bufio.NewReader(b)
	_, err = fReader.Read(buf)

	if err != nil {
		fmt.Print(err)
	}

	shoot := datashoot{}
	shoot.FileID = result.FileID
	shoot.Base64data = base64.StdEncoding.EncodeToString(buf)
	shoot.Filetype = result.Filetype
	shoot.Dataq = result.Dataq

	if templateman ==false{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(shoot)

		if err != nil{
			log.Error(fmt.Errorf("error occur in endocing json results\n"))
			log.Fatal(err.Error())
		}
	} else{
		if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
			log.Printf("unable to parse image template.")
		} else {
			data := map[string]interface{}{"Image": shoot.Base64data}
			if err = tmpl.Execute(w, data); err != nil {
				log.Printf("unable to execute template.")
			}
		}
	}
	fmt.Println(r.RemoteAddr + "Request random data from  "+ dsName + time.Now().UTC().String())


}

//func postAnswerToDataset(w http.ResponseWriter, r *http.Request){
//	_, err := fmt.Fprint(w,"ohyesgepostAnswerToDataset")
//	if err != nil{
//		log.Fatal(err.Error())
//	}
//}
//
//func getDataSetInfo(w http.ResponseWriter, r *http.Request){
//	_, err := fmt.Fprint(w,"ohyesgepostgetDataInfo")
//	if err != nil
//		log.Fatal(err.Error())
//	}
//}