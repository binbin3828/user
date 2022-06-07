/*
 * @Autor: Bobby
 * @Description: main function
 * @Date: 2022-06-06 10:32:45
 * @LastEditTime: 2022-06-07 21:47:02
 * @FilePath: \User\cmd\main.go
 */

package main

import (
	"log"
	"net/http"
	"user/pkg/dbconn"
	"user/pkg/logger"
	"user/service"

	"github.com/gorilla/mux"
)

func main() {
	logger.InitLogger()
	defer logger.SugarLogger.Sync()
	dbconn.InitMysql()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range service.AllRoutes {
		var handler http.Handler
		handler = route.HandlerFunc

		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
