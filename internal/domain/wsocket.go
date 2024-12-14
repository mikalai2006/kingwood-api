package domain

// Message struct to hold message data
type Message struct {
	Type      string      `json:"type"`
	Method    string      `json:"method"`
	Sender    string      `json:"sender"`
	Recipient string      `json:"recipient"`
	Service   string      `json:"service"`
	Content   interface{} `json:"content"`
	ID        string      `json:"id"`
}
