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
	FileID string //
	Path string
	//Filetype string
	//dataAnswer boxAnswer
	//IsFake bool
	//NumofReq int
	//NumofAns int

}

type dataset struct {
	datacoll *MongoDBCollection
	DatasetName string
	isBait bool
	numofdata int
	answerType string
}

type client struct {
	email string
	pw string
	token string
	points int
	banpoint int
	isbanned bool
	isAdmin bool
}
