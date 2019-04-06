package main

import (
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"labeltong-server/imageserver/handler"
	"net/http"
)

func main(){
	r:= mux.NewRouter()
	handler.AuthInitSubrouter(r)
	handler.DatasetInitSubrouter(r)
//	infoRouter := r.PathPrefix("/info").Subrouter()


	err := http.ListenAndServe(":19432", r)
	if err !=nil{
		log.Fatal(err.Error())
	}




}
