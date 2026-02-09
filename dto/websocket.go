package dto

// Message is the DTO for WebSocket signaling messages
type Message struct {
	Type     string      `json:"type"`
	From     string      `json:"from,omitempty"`
	To       string      `json:"to,omitempty"`
	RoomID   string      `json:"roomId,omitempty"`
	Username string      `json:"username,omitempty"`
	Payload  interface{} `json:"payload,omitempty"`
}

// SDPMessage represents an SDP offer/answer
type SDPMessage struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

// ICECandidateMessage represents an ICE candidate
type ICECandidateMessage struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

// MessageType constants for signaling
const (
	MessageTypeOffer     = "offer"
	MessageTypeAnswer    = "answer"
	MessageTypeCandidate = "candidate"
	MessageTypeJoin      = "join"
	MessageTypeLeave     = "leave"
	MessageTypeReady     = "ready"
	MessageTypeError     = "error"
)
