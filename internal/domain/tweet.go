package domain

type CreateTweetInput struct {
	Text       string `json:"text" validate:"required"`
	ImageChunk string `json:"image" validate:"base64"`
}

type TweetResponse struct {
	TweetID string `json:"tweet_id"`
	UserID  string `json:"user_id"`
	Text    string `json:"text"`
}

type IncomingNewTweetNotification struct {
	FromUserID string `json:"from_user_id"`
	ShortTitle string `json:"short_title"`
}
