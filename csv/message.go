package csv

import (
	"encoding/csv"
	"os"
	"reflect"

	"github.com/mailtomation"
)

type MessageService struct {
	csv *CSV
}

var _ mailtomation.MessageCSVService = &MessageService{}

func NewMessageService(csv *CSV) *MessageService {
	return &MessageService{csv: csv}
}

func (s *MessageService) Marshal(msg ...mailtomation.Message) ([][]string, error) {
	var collection [][]string
	for _, m := range msg {
		v := reflect.ValueOf(m)
		values := make([][]byte, v.NumField())
		strVal := make([]string, v.NumField())
		for i := range values {
			values[i] = v.Field(i).Bytes()
			strVal[i] = string(values[i])
		}
		collection = append(collection, strVal)
	}
	return collection, nil
}

func (s *MessageService) Write(data [][]string) error {
	file, err := os.OpenFile(s.csv.path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	if err := writer.WriteAll(data); err != nil {
		return err
	}
	return nil
}
