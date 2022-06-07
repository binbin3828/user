/*
 * @Autor: Bobby
 * @Description: Custom JSON time conversion tool
 * @Date: 2022-06-06 22:15:49
 * @LastEditTime: 2022-06-07 21:58:04
 * @FilePath: \User\pkg\util\JsonTime.go
 */

package util

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

// 自定义时间
type JsonTime time.Time

func (t *JsonTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	// 前端接收的时间字符串
	str := string(data)
	// 去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = JsonTime(t1)
	return err
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := "null"
	if !time.Time(t).IsZero() {
		formatted = fmt.Sprintf("\"%v\"", time.Time(t).Format("2006-01-02 15:04:05"))
	}
	return []byte(formatted), nil
}

func (t JsonTime) Value() (driver.Value, error) {
	// 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006-01-02 15:04:05"), nil
}

func (t *JsonTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = JsonTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *JsonTime) String() string {
	return time.Time(*t).String()
}
