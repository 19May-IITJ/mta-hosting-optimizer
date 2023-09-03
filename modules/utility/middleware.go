package utility

type Message struct {
	Hostname string `json:"hostname"`
	Active   int    `json:"active"`
	// Passive  int    `json:"passive"`
}

func NewMessage() *Message {
	return &Message{}
}
