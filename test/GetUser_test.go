/*
 * @Autor: Bobby
 * @Description: unit test GetUser
 * @Date: 2022-06-08 11:45:03
 * @LastEditTime: 2022-06-08 14:57:07
 * @FilePath: \user\test\GetUser_test.go
 */
package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// run test with the following command:
// go test -v .\GetUser_test.go
func TestGetUser_Run(t *testing.T) {
	url := "http://127.0.0.1:8080/user/2"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
	fmt.Println(string(body))
}
