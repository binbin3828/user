package main

import (
	"log"
	"net/http"
	"user/dao"
	"user/pkg/dbconn"
	"user/pkg/logger"
	"user/service"
)

func main() {
	zapLog := logger.NewZapLogger()
	defer zapLog.Sync()

	db, err := dbconn.NewMysql(zapLog)
	if err != nil {
		log.Fatalf("mysql init failed: %v", err)
	}

	userDao := dao.NewUserDao(db, zapLog)
	friendsDao := dao.NewFriendsDao(db, zapLog)

	svc := service.NewService(zapLog, userDao, friendsDao)
	router := service.NewRouter(svc)

	log.Fatal(http.ListenAndServe(":8080", router))
}
