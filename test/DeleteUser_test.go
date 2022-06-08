/*
 * @Autor: Bobby
 * @Description:unit test DeleteUser
 * @Date: 2022-06-08 14:45:13
 * @LastEditTime: 2022-06-08 14:56:55
 * @FilePath: \user\test\DeleteUser_test.go
 */
package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// run test with the following command:
// go test -v .\DeleteUser_test.go
func TestDeleteUser_Run(t *testing.T) {
	url := "http://127.0.0.1:8080/user/3"
	req, _ := http.NewRequest("DELETE", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.Status)
	fmt.Println(string(body))
}
