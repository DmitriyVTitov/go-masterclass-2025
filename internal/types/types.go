package types

type Review struct {
	ObjectID int    `json:"object_id"`
	Text     string `json:"text,omitempty"`
}
