package mailtomation

type MessageService interface {
	ReadUnreadMessages(user string) ([]Message, error)
}

type MessageCSVService interface {
	Marshal(m ...Message) ([][]string, error)
}

type field []byte

type Message struct {
	Sender field
	Plain  field
	HTML   field
}

func (f field) String() string {
	return string(f)
}

const (
	NoUnreadMessages = Error("No unread messages")
)
