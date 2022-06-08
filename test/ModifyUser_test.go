/*
 * @Autor: Bobby
 * @Description: unit test ModifyUser
 * @Date: 2022-06-08 14:50:28
 * @LastEditTime: 2022-06-08 14:57:28
 * @FilePath: \user\test\ModifyUser_test.go
 */
package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"user/model"
)

// run test with the following command:
// go test -v .\ModifyUser_test.go
func TestModifyUser_Run(t *testing.T) {
	url := "http://127.0.0.1:8080/user"

	user := model.User{
		Id:      10,
		Name:    "bobby2",
		Address: "guangzhou",
	}
	sbyte, _ := json.Marshal(user)
	reader := strings.NewReader(string(sbyte))
	req, _ := http.NewRequest("PUT", url, reader)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
	fmt.Println(string(body))
}
