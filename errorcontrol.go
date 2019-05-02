package main

import "log"

func customfatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func custompanic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

