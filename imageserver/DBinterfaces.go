package main

type DBinterface interface {
	InsertData(data interface{} ) (err error)
	RemoveData(dataId string) (err error)
	ReadData(dataId string) (data interface{}, err error)
}
