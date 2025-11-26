package models

type SessionResponse struct {
	SessionId string `json:"sessionId"`
}

type UpsertProductVoteRequest struct {
	ProductId   string  `json:"productId"`
	ProductName string  `json:"productName"`
	MachineId   *string `json:"machineId"`
	SessionId   string  `json:"-"`
	Liked       bool    `json:"liked"`
}

type UpsertProductVoteResponse struct {
	ProductId string `json:"productId"`
	Message   string `json:"message"`
}
