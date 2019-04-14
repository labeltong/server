package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main(){
	r:= mux.NewRouter()
	AuthInitSubrouter(r)
	DatasetInitSubrouter(r)
	//	infoRouter := r.PathPrefix("/info").Subrouter()
	r.HandleFunc("/",func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		_, err := w.Write([]byte("Hello Labeltong user"))
		if err != nil{
			log.Error(fmt.Errorf("Log error in /" ))
			log.Fatal(err.Error())
		}

	})

	err := http.ListenAndServe(":19432", r)
	if err !=nil{
		log.Fatal(err.Error())
	}


}
