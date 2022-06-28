/*
 * @Autor: Bobby
 * @Description: test add friends
 * @Date: 2022-06-09 16:27:04
 * @LastEditTime: 2022-06-09 16:42:21
 * @FilePath: \user\test\AddFriends_test.go
 */
package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestAddFriend_Run(t *testing.T) {
	url := "http://127.0.0.1:8080/friends"
	friends := make(map[string]interface{})
	friends["uid"] = 9
	friends["fri"] = 10
	sbyte, _ := json.Marshal(friends)
	reader := strings.NewReader(string(sbyte))
	req, _ := http.NewRequest("POST", url, reader)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
	fmt.Println(string(body))
}
