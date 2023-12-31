package domain

import "mime/multipart"

type CreateCommentInput struct {
	TweetID    string `json:"tweet_id" validate:"uuid4"`
	Text       string `json:"text" validate:"required"`
	ImageChunk string `json:"image" validate:"base64"`
}

type CreateCommentRequest struct {
	UserID  string                `json:"user_id"`
	TweetID string                `json:"tweet_id"`
	Text    string                `json:"text"`
	Image   *multipart.FileHeader `json:"image"`
}

type CommentResponse struct {
	CommentID string `json:"comment_id"`
	UserID    string `json:"user_id,omitempty"`
	TweetID   string `json:"tweet_id"`
	Text      string `json:"text"`
	ImageURL  string `json:"image,omitempty"`
}

type UpdateCommentRequest struct {
	TweetID   string
	UserID    string
	Text      string
	Image     *multipart.FileHeader
	CommentID string
}
