package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	"time"
)

var (
	t        *string
	filePath *string
)

func main() {
	t = pflag.StringP("type", "t", "", "operation type")
	filePath = pflag.StringP("file_path", "f", "", "specify a path to upload")
	pflag.Parse()

	log.Print(*filePath)
	viper.BindEnv("SecretKey", "SecretKey")
	viper.BindEnv("SecretId", "SecretId")
	viper.BindEnv("MaxSize", "MaxSize")
	viper.BindEnv("ImageBucket", "ImageBucket")
	viper.BindEnv("AppId", "AppId")

	switch *t {
	case "upload":
		RunUpload()
	case "list":
		RunList()
	}
}

func addAuth(request *http.Request) {
	authTime := cos.NewAuthTime(time.Minute)
	cos.AddAuthorizationHeader(viper.GetString("SecretId"), viper.GetString("SecretKey"),
		"", request, authTime)
}
