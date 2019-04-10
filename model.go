package main

type boxAnswer struct {
	x1 int
	y1 int
	x2 int
	y2 int
}
type classifyAnswer struct {
	ans int
}

type datatolabel struct {
	FileNum int `json:"file_num"`
	FileID string `json:"file_id"`
	Path   string `json:"path"`
	//Filetype string
	//dataAnswer boxAnswer
	//IsFake bool
	//NumofReq int
	//NumofAns int

}

type datashoot struct{
	FileID string `json:"file_id"`
	Base64data string `json:"base_64_data"`
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
