package main

import "google.golang.org/api/gmail/v1"

type Config struct {
	Scopes []string
}

func NewConfig() *Config {
	return &Config{
		Scopes: []string{
			gmail.GmailReadonlyScope,
			gmail.GmailLabelsScope,
			gmail.GmailModifyScope,
		},
	}
}
