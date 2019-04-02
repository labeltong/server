package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"testing"
)

type testdata struct{
	id string
	datapath string
}

func TestMongoDBCollection_InsertData(t *testing.T) {
	d,e := InitDatabase("test","labeltong")
	yesman := datatolabel{
		FileID:"1",
		Path:"dsfasdfsfsad",
	}
	if e !=nil{
		t.Error("fail initserver")
		fmt.Println(e.Error())
		log.Fatal(e.Error())
	}
	e = d.InsertData(yesman)
	if e !=nil{
		t.Error("fail intesrt")
		fmt.Println(e.Error())
		log.Fatal(e.Error())
	}

}

func TestMongoDBCollection_ReadData(t *testing.T) {
	d,e := InitDatabase("test","labeltong")
	if e !=nil{
		t.Error("fail initserver")
		fmt.Println(e.Error())
		log.Fatal(e.Error())
	}

	res,err := d.ReadData("1")
	if err != nil{
		t.Error("fail read")
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}
	fmt.Println(res)

}

func TestMongoDBCollection_RemoveData(t *testing.T) {
	d,e := InitDatabase("test","labeltong")
	yesman := datatolabel{
		FileID:"56785678",
		Path:"myresult5678",
	}
	if e !=nil{
		t.Error("fail initserver")
		fmt.Println(e.Error())
		log.Fatal(e.Error())
	}
	e = d.RemoveData(yesman.FileID)
	if e !=nil{
		t.Error("fail delete")
		fmt.Println(e.Error())
		log.Fatal(e.Error())
	}

}
