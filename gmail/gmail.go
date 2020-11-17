package gmail

import (
	"net/http"

	"google.golang.org/api/gmail/v1"
)

type Gmail struct {
	g *gmail.Service
}

func New(c *http.Client) (*Gmail, error) {
	s, err := gmail.New(c)
	if err != nil {
		return nil, err
	}
	return &Gmail{g: s}, nil
}
