package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"
)

func RunUpload() {
	upload(*filePath)
}

func upload(filePath string) {
	fileBytes, baseName := readFile(filePath)
	if fileBytes == nil {
		return
	}
	req := constructObjectReq(fileBytes, baseName)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	dumpResponse, err := httputil.DumpResponse(rsp, true)
	if err == nil {
		log.Println(string(dumpResponse))
	}
	if rsp.StatusCode != 200 {
		os.Stdout.Write([]byte("upload failed"))
	} else {
		os.Stdout.Write([]byte(req.URL.String()))
	}
}

func constructObjectReq(fileBytes []byte, fileName string) *http.Request {
	contentMd5 := calMd5(fileBytes)

	u := viper.GetString("ImageBucket") + "/" + strconv.Itoa(int(time.Now().Unix())) + "_" + fileName
	request, err := http.NewRequest("PUT", u, bytes.NewReader(fileBytes))
	if err != nil {
		log.Fatal(err)
	}
	//request.Header.Add("Content-Type", "image/jpeg")
	request.Header.Add("Content-MD5", contentMd5)

	addAuth(request)

	dumpRequest, err := httputil.DumpRequest(request, true)
	if err == nil {
		log.Println(string(dumpRequest))
	}
	return request
}

func calMd5(fileBytes []byte) string {
	hash := md5.New()
	_, err := hash.Write(fileBytes)
	if err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)
	md5Content := base64.StdEncoding.EncodeToString(sum)
	return md5Content
}

func readFile(filePath string) ([]byte, string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileBytes := make([]byte, fileInfo.Size())

	if len(fileBytes) > viper.GetInt("MaxSize")*1000 { // MaxSize单位：1kb
		os.Stdout.Write([]byte(fmt.Sprintf("Max upload size is set to %dKB", viper.GetInt("MaxSize"))))
		return nil, ""
	}
	n, err := file.Read(fileBytes)
	if err != nil || n == 0 {
		log.Fatal("err: ", err, "\n write n = ", n)
	}
	return fileBytes, fileInfo.Name()
}
