/*
 * @Autor: Bobby
 * @Description: routes for API
 * @Date: 2022-06-06 11:01:16
 * @LastEditTime: 2022-06-07 22:01:41
 * @FilePath: \User\service\Routes.go
 */
package service

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var AllRoutes = Routes{
	Route{Name: "User", Method: "GET", Pattern: "/user/{uid}", HandlerFunc: GetUser},
	Route{Name: "User", Method: "POST", Pattern: "/user", HandlerFunc: CreateUser},
	Route{Name: "User", Method: "DELETE", Pattern: "/user/{uid}", HandlerFunc: DeleteUser},
	Route{Name: "User", Method: "PUT", Pattern: "/user", HandlerFunc: ModifyUser},
}
