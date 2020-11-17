package gmail

import (
	"encoding/base64"

	"github.com/mailtomation"
)

type MessageService struct {
	g *Gmail
}

var _ mailtomation.MessageService = &MessageService{}

func NewMessageService(g *Gmail) *MessageService {
	return &MessageService{g: g}
}

func (s *MessageService) ReadUnreadMessages(user string) (
	[]mailtomation.Message, error) {
	var collection []mailtomation.Message
	messages := s.g.g.Users.Messages
	messageList, err := messages.List(user).Q(
		"is:unread").Do()
	if err != nil {
		return nil, err
	}
	if len(messageList.Messages) == 0 {
		return nil, mailtomation.NoUnreadMessages
	}
	for _, m := range messageList.Messages {
		var newMessage mailtomation.Message
		messageContent, err := messages.Get(user, m.Id).Do()
		if err != nil {
			return nil, err
		}
		for _, part := range messageContent.Payload.Parts {
			switch part.MimeType {
			case "text/plain":
				if newMessage.Plain, err = decodeBase64URL(part.Body.Data); err != nil {
					return nil, err
				}
			case "text/html":
				if newMessage.HTML, err = decodeBase64URL(part.Body.Data); err != nil {
					return nil, err
				}
			default:
				continue
			}
		}
		// if _, err := messages.Modify(user, m.Id, &gmail.ModifyMessageRequest{
		// 	AddLabelIds:    []string{},
		// 	RemoveLabelIds: []string{"UNREAD"},
		// }).Do(); err != nil {
		// 	return nil, err
		// }
		collection = append(collection, newMessage)
	}
	return collection, nil
}

func decodeBase64URL(str string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(str)
}
