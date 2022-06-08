/*
 * @Autor: Bobby
 * @Description: unit test CreateUser
 * @Date: 2022-06-08 14:33:13
 * @LastEditTime: 2022-06-08 14:57:01
 * @FilePath: \user\test\CreateUser_test.go
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
// go test -v .\CreateUser_test.go
func TestCreateUser_Run(t *testing.T) {
	url := "http://127.0.0.1:8080/user"
	user := model.User{
		Name:        "bobby1",
		Dob:         "1990-01-10",
		Address:     "shenzhen",
		Description: "coder",
	}
	sbyte, _ := json.Marshal(user)
	reader := strings.NewReader(string(sbyte))
	req, _ := http.NewRequest("POST", url, reader)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
	fmt.Println(string(body))
}
