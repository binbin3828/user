/*
 * @Autor: Bobby
 * @Description: unit test for friends list
 * @Date: 2022-06-09 16:32:31
 * @LastEditTime: 2022-06-09 18:39:46
 * @FilePath: \user\test\GetFriendsList_test.go
 */
package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetFriendsList_Run(t *testing.T) {
	url := "http://127.0.0.1:8080/friends/10"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
	fmt.Println(string(body))
}
