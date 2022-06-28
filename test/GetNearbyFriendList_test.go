/*
 * @Autor: Bobby
 * @Description: unit test find nearby friends
 * @Date: 2022-06-09 18:37:47
 * @LastEditTime: 2022-06-09 18:37:50
 * @FilePath: \user\test\GetNearbyFriend.go
 */

package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetNearbyFriendsList_Run(t *testing.T) {
	url := "http://127.0.0.1:8080/nearbyfriends/10"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
	fmt.Println(string(body))
}
