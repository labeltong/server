package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)
/*
/dataset/list [GET] Login function, Oauth from client required
/list/{id}/get [GET] Logout function, Oauth from client required
/list/{id}/ans [POST] Check if user is authenticated
/list/{id}/info [GET] Check if user is authenticated

*/
func DatasetInitSubrouter(r *mux.Router)  {
	ret := r.PathPrefix("/dataset").Subrouter()

	ret.HandleFunc("/list", listAlldatasets).Methods("GET")
	ret.HandleFunc("/list/{id}/get", getDataFromDataset).Methods("GET")
	ret.HandleFunc("/list/{id}/ans",postAnswerToDataset).Methods("POST")
	ret.HandleFunc("/list/{id}/info",getDataInfo).Methods("GET")


}


func listAlldatasets(w http.ResponseWriter, r *http.Request){

}


func getDataFromDataset(w http.ResponseWriter, r *http.Request){

}

func postAnswerToDataset(w http.ResponseWriter, r *http.Request){

}

func getDataInfo(w http.ResponseWriter, r *http.Request){

}