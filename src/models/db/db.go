package db

import (
	"time"
)

type App struct {
	Id   int
	Name string
}

type OnlineHistory struct {
	AppId    int
	Count    int
	Datetime time.Time
}
