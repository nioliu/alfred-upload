package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	"time"
)

type Result struct {
	Items []*Item `json:"items"`
}

type Item struct {
	Uid          string            `json:"uid"`
	Type         string            `json:"type"`
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle"`
	Arg          string            `json:"arg"` // 相当于query
	Autocomplete string            `json:"autocomplete"`
	Icon         *Icon             `json:"icon"`
	Variables    map[string]string `json:"variables"`
	QuickLookUrl string            `json:"quicklookurl"`
	Text         *Text             `json:"text"` // Defines the text the user will get when copying the selected result row with ⌘C or displaying large type with ⌘L.
}

type Text struct {
	Copy      string `json:"copy"`
	Largetype string `json:"largetype"`
}

type Icon struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

var (
	t        *string
	filePath *string
)

func main() {
	t = pflag.StringP("type", "t", "", "operation type")
	filePath = pflag.StringP("file_path", "f", "", "specify a path to upload")
	pflag.Parse()

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	log.Print(viper.AllSettings())
	log.Print(*filePath)

	switch *t {
	case "upload":
		RunUpload()
	case "list":
		RunList()
	}
}

func PrintError(err error) {
	errorPrint := GetErrorPrint(err)
	fmt.Println(errorPrint)
	log.Fatal(errorPrint)
}

func GetErrorPrint(err error) string {
	return fmt.Sprintf(`{"items":[{"uid":"error","type":"text","title":"Something Error","subtitle":"%s"}]}`, err.Error())
}

func addAuth(request *http.Request) {
	authTime := cos.NewAuthTime(time.Minute)
	cos.AddAuthorizationHeader(viper.GetString("SecretId"), viper.GetString("SecretKey"),
		"", request, authTime)
}
