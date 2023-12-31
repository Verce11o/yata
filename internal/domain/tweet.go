package domain

import "mime/multipart"

type CreateTweetInput struct {
	Text       string `json:"text" validate:"required"`
	ImageChunk string `json:"image" validate:"base64"`
}

type TweetResponse struct {
	TweetID string `json:"tweet_id"`
	UserID  string `json:"user_id"`
	Text    string `json:"text"`
}

type CreateTweetRequest struct {
	UserID string
	Text   string
	Image  *multipart.FileHeader
}

type UpdateTweetRequest struct {
	UserID  string
	TweetID string
	Text    string
	Image   *multipart.FileHeader
}
