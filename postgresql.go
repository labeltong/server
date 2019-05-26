package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/gommon/log"
	"time"
)


type postgredb struct{
	DB *gorm.DB
	db_type string
	host string
	port int
	user string
	dbname string
	password string
	maxidleconn int
	maxopenconn int
	maxconnhour time.Duration


}


func (p *postgredb) Connect() (err error){

	p.db_type = "postgres"

	p.host = DB_host
	p.port = DB_port
	p.user = DB_user
	p.dbname = DB_dbname
	p.password = DB_password
	p.maxidleconn = DB_maxidlevalue
	p.maxopenconn = DB_maxopenvalue
	p.maxconnhour = time.Hour

	connstr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s  password=%s sslmode=disable",p.host, p.port, p.user, p.dbname, p.password)
	p.DB, err = gorm.Open(p.db_type, connstr)

	if err != nil {
		log.Panic(err)
		return err
	}

	p.DB.DB().SetMaxIdleConns(p.maxidleconn)
	p.DB.DB().SetMaxOpenConns(p.maxopenconn)
	p.DB.DB().SetConnMaxLifetime(p.maxconnhour)
	p.DB.LogMode(true)


	return nil

}

