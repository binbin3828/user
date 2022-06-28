/*
 * @Autor: Bobby
 * @Description: unit test ModifyUser
 * @Date: 2022-06-08 14:50:28
 * @LastEditTime: 2022-06-09 17:59:33
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
)

// run test with the following command:
// go test -v .\ModifyUser_test.go
func TestModifyUser_Run(t *testing.T) {
	url := "http://127.0.0.1:8080/user"

	user := make(map[string]interface{})
	user["id"] = 4
	// user["name"] = "bobby4"
	// user["address"] = "shenzhen"
	user["latitude"] = 39.911987
	user["longitude"] = 116.414311
	sbyte, _ := json.Marshal(user)
	reader := strings.NewReader(string(sbyte))
	req, _ := http.NewRequest("PUT", url, reader)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
	fmt.Println(string(body))
}
