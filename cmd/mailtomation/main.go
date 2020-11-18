package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mailtomation/csv"
	"github.com/mailtomation/gmail"
	"github.com/mailtomation/oauth2"
)

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config := NewConfig()
	oauth2Service, err := oauth2.New(b, config.Scopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	s, err := gmail.New(oauth2Service.GetClient())
	messageService := gmail.NewMessageService(s)
	msgs, err := messageService.ReadUnreadMessages("me")
	if err != nil {
		fmt.Println(err)
	}
	CSVService := csv.New("sample.csv")
	messageCSVService := csv.NewMessageService(CSVService)
	data, _ := messageCSVService.Marshal(msgs...)
	if err := messageCSVService.Write(data); err != nil {
		log.Fatal(err)
	}
}