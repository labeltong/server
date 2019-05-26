package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email         string
	Token         string
	Name          string
	PhoneNum      string
	Points        int
	IsBanned      bool
	BanPoint      int
	IsAdmin       bool
	LastLoginDate time.Time
	AnswerList    []Answers `gorm:"foreignkey:UserId"`
	PointUsageList []PointUsage `gorm:"foreignkey:UserId"`
}


type PointUsage struct{
	gorm.Model
	UserId uint
	ItemId uint

}



type Items struct{
	gorm.Model
	Name string
	Price int
	Thumbnails string
}

type Datasets struct {
	gorm.Model
	DatasetName        string
	DatasetDescription string
	Owner              int
	IsFake             bool
	NumberOfDataset    int
	DatasetThumbnail   string
	Answertype         string
	DataList []Datas `gorm:"foreignkey:DatasetId"`
}

type Datas struct {
	gorm.Model
	DatasetId         uint
	IsFake            bool
	RequiredNumAnswer int
	DataPath          string
	AnswerType        string
	Question          string
	TagList []Tags `gorm:"foreignkey:DataId"`
}

type TagList struct {
	gorm.Model
	TagId uint `gorm:"id"`
	TagName string
	TagDescription string
	TagThumbnail string
}

type Tags struct{
	gorm.Model
	DataID uint
	TagID uint
}


type Answers struct {
	gorm.Model
	UserId     uint
	DataId     uint
	IsValid    bool
	AnswerData string
	AnswerTime time.Time
}

type AllAnswers struct {
	gorm.Model
	UserId     uint
	DataId uint
	IsValid    bool
	IsBait bool
	AnswerData string
	AnswerTime time.Time

}

type AnswersJSON struct {
	Email     string `json:"email"`
	DataId     uint `json:"data_id"`
	AnswerData string `json:"answer_data"`
}

type UserInfoClaim struct {
	Email         string    `json:"email"`
	Token         string    `json:"token"`
	Name          string    `json:"name"`
	PhoneNum      string    `json:"phone_num"`
	Points        int       `json:"points"`
	IsBanned      bool      `json:"is_banned"`
	LastLoginDate time.Time `json:"last_login_date"`
	RegisterDate  time.Time `json:"register_date"`
	jwt.StandardClaims
}

type MakeUser struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	PhoneNum string `json:"phone_num"`
}
