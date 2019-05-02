package main

import "github.com/dgrijalva/jwt-go"

type boxAnswer struct {
	x1 int
	y1 int
	x2 int
	y2 int
}
type classifyAnswer struct {
	questions []string
	qlength   int
	//ans int
}

type sentimentAnswer struct {
	questions []string
	qlength   int
	//ans int
}

type answertype struct {
	classify  classifyAnswer
	sentiment sentimentAnswer
	box       boxAnswer
}

type datatolabel struct {
	FileNum  int    `json:"file_num"`
	FileID   string `json:"file_id"`
	Path     string `json:"path"`
	Filetype string // classify, sentiment, box
	Dataq    []string
	//IsFake bool
	//NumofReq int
	//NumofAns int

}

type datashoot struct {
	FileID     string `json:"file_id"`
	Base64data string `json:"base_64_data"`
	Filetype   string // classify, sentiment, box
	Dataq      []string
}

type dataset struct {
	Datasetname   string `json:"datasetname"`
	Datasetdbname string `json:"datasetdbname"`
	Isbait        bool   `json:"isbait"`
	Numofdata     int    `json:"numofdata"`
	Answertype    string `json:"answertype"`
}

type clientdata struct {
	Email      string   `json:"email"`
	Token      string   `json:"token"`
	Points     int      `json:"points"`
	Banpoint   int      `json:"banpoint"`
	Isbanned   bool     `json:"isbanned"`
	IsAdmin    bool     `json:"is_admin"`
	Pointusage []string `json:"pointusage"`
}

type clienttouser struct {
	Email      string   `json:"email"`
	Token      string   `json:"token"`
	Points     int      `json:"points"`
	Isbanned   bool     `json:"isbanned"`
	Pointusage []string `json:"pointusage"`
}
type user struct {
	Email      string   `json:"email"`
	Token      string   `json:"token"`

}
type Claims struct {
	Email string `json:"email"`
	Points     int      `json:"points"`
	Isbanned   bool     `json:"isbanned"`
	Pointusage []string `json:"pointusage"`
	jwt.StandardClaims
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}