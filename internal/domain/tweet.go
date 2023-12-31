package domain

import "time"

type CreateTweetInput struct {
	Text       string `json:"text" validate:"required"`
	ImageChunk string `json:"image" validate:"base64"`
}

type TweetResponse struct {
	TweetID   string    `json:"tweet_id"`
	UserID    string    `json:"user_id,omitempty"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTweetRequest struct {
	UserID string
	Text   string
	Image  *Image
}

type UpdateTweetRequest struct {
	UserID  string
	TweetID string
	Text    string
	Image   *Image
}
