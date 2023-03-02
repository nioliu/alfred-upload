package main

import (
	"github.com/spf13/viper"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	log.Print(viper.AllSettings())

	m.Run()
}
func TestRun(t *testing.T) {
	RunUpload()
}

func TestFi(t *testing.T) {
	readFile("./sort.png")
}

func TestLis(t *testing.T) {
	RunList()
}
