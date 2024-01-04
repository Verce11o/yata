package domain

import "time"

type CreateCommentInput struct {
	TweetID    string `json:"tweet_id" validate:"uuid4"`
	Text       string `json:"text" validate:"required"`
	ImageChunk string `json:"image" validate:"base64"`
}

type CreateCommentRequest struct {
	UserID  string
	TweetID string
	Text    string
	Image   *Image
}

type CommentResponse struct {
	CommentID string    `json:"comment_id"`
	UserID    string    `json:"user_id,omitempty"`
	TweetID   string    `json:"tweet_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateCommentRequest struct {
	TweetID   string
	UserID    string
	Text      string
	Image     *Image
	CommentID string
}
