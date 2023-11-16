package items

import "strings"

type Item struct {
	ID           string `json:"id"`
	CollectionId string `json:"collection_id"`
	UserId       string `json:"-"`
	URL          string `json:"url"`
	Title        string `json:"title"`
	Description  string `json:"desc"`
	CreatedAt    int64  `json:"created_at"`
}

func (i *Item) Validate() error {
	i.URL = strings.TrimSpace(i.URL)
	i.Title = strings.TrimSpace(i.Title)
	i.Description = strings.TrimSpace(i.Description)

	if i.UserId == "" {
		return ErrMissingUserId
	}

	if i.URL == "" {
		return ErrMissingURL
	}

	if i.Title == "" {
		return ErrMissingTitle
	}

	return nil
}
