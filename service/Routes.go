/*
 * @Autor: Bobby
 * @Description: routes for API
 * @Date: 2022-06-06 11:01:16
 * @LastEditTime: 2022-06-16 13:45:46
 * @FilePath: \user\service\Routes.go
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
	//user
	Route{Name: "User", Method: "GET", Pattern: "/user/{uid}", HandlerFunc: responseHandler(GetUser)},
	Route{Name: "User", Method: "POST", Pattern: "/user", HandlerFunc: responseHandler(CreateUser)},
	Route{Name: "User", Method: "DELETE", Pattern: "/user/{uid}", HandlerFunc: responseHandler(DeleteUser)},
	Route{Name: "User", Method: "PUT", Pattern: "/user", HandlerFunc: responseHandler(ModifyUser)},

	//friends
	Route{Name: "Friends", Method: "POST", Pattern: "/friends", HandlerFunc: responseHandler(AddFriend)},
	Route{Name: "Friends", Method: "GET", Pattern: "/friends/{uid}", HandlerFunc: responseHandler(GetFriendsList)},
	Route{Name: "NearbyFriends", Method: "GET", Pattern: "/nearbyfriends/{uid}", HandlerFunc: responseHandler(GetNearbyFriend)},
}
