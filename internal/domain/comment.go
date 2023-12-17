package domain

type CreateCommentInput struct {
	TweetID    string `json:"tweet_id" validate:"uuid4"`
	Text       string `json:"text" validate:"required"`
	ImageChunk string `json:"image" validate:"base64"`
}

type CommentResponse struct {
	CommentID string `json:"comment_id"`
	UserID    string `json:"user_id,omitempty"`
	TweetID   string `json:"tweet_id"`
	Text      string `json:"text"`
	ImageURL  string `json:"image,omitempty"`
}
