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

// AllRoutes 根据注入的 Service 构建所有路由
func AllRoutes(s *Service) Routes {
	return Routes{
		Route{Name: "User", Method: "GET", Pattern: "/user/{uid}", HandlerFunc: s.responseHandler(s.GetUser)},
		Route{Name: "User", Method: "POST", Pattern: "/user", HandlerFunc: s.responseHandler(s.CreateUser)},
		Route{Name: "User", Method: "DELETE", Pattern: "/user/{uid}", HandlerFunc: s.responseHandler(s.DeleteUser)},
		Route{Name: "User", Method: "PUT", Pattern: "/user", HandlerFunc: s.responseHandler(s.ModifyUser)},

		Route{Name: "Friends", Method: "POST", Pattern: "/friends", HandlerFunc: s.responseHandler(s.AddFriend)},
		Route{Name: "Friends", Method: "GET", Pattern: "/friends/{uid}", HandlerFunc: s.responseHandler(s.GetFriendsList)},
		Route{Name: "NearbyFriends", Method: "GET", Pattern: "/nearbyfriends/{uid}", HandlerFunc: s.responseHandler(s.GetNearbyFriend)},
	}
}
