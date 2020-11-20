package gmail

import (
	"encoding/base64"

	"github.com/nikunicke/mailtomation"
	"google.golang.org/api/gmail/v1"
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
	msgService := s.g.g.Users.Messages
	msgList, err := msgService.List(user).Q("is:unread").Do()
	if err != nil {
		return nil, err
	} else if len(msgList.Messages) == 0 {
		return nil, mailtomation.NoUnreadMessages
	}
	return readMessages(user, msgService, msgList)
}

func readMessages(user string,
	s *gmail.UsersMessagesService,
	list *gmail.ListMessagesResponse) ([]mailtomation.Message, error) {
	var collection []mailtomation.Message
	for _, m := range list.Messages {
		message, err := s.Get(user, m.Id).Do()
		if err != nil {
			return nil, err
		}
		messageData, err := readMessage(message)
		if err != nil {
			return nil, err
		}
		collection = append(collection, *messageData)
		if err := markAsRead(user, m.Id, s); err != nil {
			return nil, err
		}
	}
	return collection, nil
}

func readMessage(m *gmail.Message) (*mailtomation.Message, error) {
	var message mailtomation.Message
	var err error
	for _, p := range m.Payload.Parts {
		switch p.MimeType {
		case "text/plain":
			if message.Plain, err = decodeBase64URL(p.Body.Data); err != nil {
				return nil, err
			}
		case "text/html":
			if message.HTML, err = decodeBase64URL(p.Body.Data); err != nil {
				return nil, err
			}
		default:
			continue
		}
	}
	message.Sender, err = recievedFrom(m)
	if err != nil {
		return &message, err
	}
	return &message, nil
}

func markAsRead(user string, id string,
	s *gmail.UsersMessagesService) error {
	if _, err := s.Modify(user, id, &gmail.ModifyMessageRequest{
		AddLabelIds:    []string{},
		RemoveLabelIds: []string{"UNREAD"},
	}).Do(); err != nil {
		return err
	}
	return nil
}

func recievedFrom(m *gmail.Message) ([]byte, error) {
	for _, h := range m.Payload.Headers {
		switch h.Name {
		case "From":
			return []byte(h.Value), nil
		default:
			continue
		}
	}
	return nil, headerNotFound("From")
}

func decodeBase64URL(str string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(str)
}

func headerNotFound(header string) mailtomation.Error {
	return mailtomation.Error("Header not found: '" + header + "'")
}
