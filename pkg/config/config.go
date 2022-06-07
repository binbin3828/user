/*
 * @Autor: Bobby
 * @Description: function to parsing .yaml file
 * @Date: 2022-06-06 15:45:11
 * @LastEditTime: 2022-06-07 21:52:10
 * @FilePath: \User\pkg\config\config.go
 */
package config

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	viper2 "github.com/spf13/viper"
)

func Get(fileKey string) interface{} {
	index := strings.Index(fileKey, ".")
	fileName := fileKey[0:index]
	key := fileKey[index+1:]

	viper := viper2.New()
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../pkg/config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	return viper.Get(key)
}
