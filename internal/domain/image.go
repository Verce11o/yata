package domain

type Image struct {
	ContentType string
	Chunk       []byte
	ImageName   string
}
