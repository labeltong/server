#!/usr/bin/env bash
# For build server
## auth
go build -v auth.go config.go model.go postgresql.go utils.go
## dataset
go build -v dataset.go config.go model.go postgresql.go utils.go

## answer
go build -v answer.go config.go model.go postgresql.go utils.go

## main
go build -v main.go config.go model.go postgresql.go utils.go

#

# execute answer
sh -c './auth :13201  & ./auth :13202 & ./auth :13203 ' &


# execute dataset


sh -c './dataset :13204 & ./dataset :13205 & ./dataset :13206  ' &


## execute answer

sh -c './answer :13207 &./answer :13208 & ./answer :13209  ' &


sh -c './main :13230' &

## execute main
