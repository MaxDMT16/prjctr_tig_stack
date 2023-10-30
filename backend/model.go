package backend

type Message struct {
	ID  string `json:"id"`
	Data string `json:"data"`
	Sender string `json:"sender"`
}
