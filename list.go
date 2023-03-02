package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func RunList() {
	standardRun()
	//selfRun()
}

func standardRun() {
	u, _ := url.Parse(viper.GetString("ImageBucket"))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  viper.GetString("SecretId"),
			SecretKey: viper.GetString("SecretKey"),
		},
	})
	res, _, err := client.Bucket.Get(context.Background(), nil)
	if err != nil {
		PrintError(err)
	}
	results := &Result{Items: make([]*Item, 0, len(res.Contents))}

	if res.Contents != nil {
		for _, content := range res.Contents {
			results.Items = append(results.Items, &Item{
				Uid:          content.ETag,
				Type:         "url",
				Title:        content.Key,
				Subtitle:     content.StorageClass,
				Arg:          viper.GetString("ImageBucket") + "/" + content.Key,
				Autocomplete: "",
				Icon:         nil,
				Variables: map[string]string{
					"file_path": viper.GetString("ImageBucket") + "/" + content.Key,
				},
				QuickLookUrl: viper.GetString("ImageBucket") + "/" + content.Key,
				Text: &Text{
					Copy:      "![](" + viper.GetString("ImageBucket") + "/" + content.Key + ")",
					Largetype: viper.GetString("ImageBucket") + "/" + content.Key,
				},
			})
		}
	}

	bytes, err := json.Marshal(results)
	if err != nil {
		PrintError(err)
	}
	fmt.Print(string(bytes))
}

// 废弃
func selfRun() {
	req := constructListReq()
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		PrintError(err)
	}
	respBytes, err := httputil.DumpResponse(rsp, true)
	if err == nil {
		println(string(respBytes))
	}
}

func constructListReq() *http.Request {
	req, err := http.NewRequest("GET", viper.GetString("ImageBucket"), nil)
	if err != nil {
		PrintError(err)
	}
	addAuth(req)

	dumpRequest, err := httputil.DumpRequest(req, true)
	if err == nil {
		println(string(dumpRequest))
	}

	return req
}
