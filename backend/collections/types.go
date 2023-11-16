package collections

import (
	"errors"
	"strings"
)

type Collection struct {
	Id         string `json:"id"`
	UserId     string `json:"-"`
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
	Archived   bool   `json:"archived"`
	ArchivedAt int64  `json:"archived_at"`
}

func (c *Collection) isValid() error {
	c.Name = strings.TrimSpace(c.Name)

	if c.Id == "" {
		return errors.New("missing collection id")
	}

	if c.UserId == "" {
		return errors.New("missing user id")
	}

	if c.Name == "" {
		return errors.New("missing name")
	}

	return nil
}
