package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"net/url"
)

func MakeReverseProxy(rplist []string, rpnum *int) (n int, t []*middleware.ProxyTarget) {
	var urltmp *url.URL
	tm := []*middleware.ProxyTarget{}
	var err error
	n = 0
	for _, s := range rplist {
		urltmp, err = url.Parse(s)
		if err != nil {
			log.Panic(err)
		}
		tm = append(tm, &middleware.ProxyTarget{
			Name: urltmp.String(),
			URL:  urltmp,
		})
		n++
	}
	return n, tm
}

func Custom_fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Custom_panic(err error) {
	if err != nil {
		log.Panic(err)
	}
}
