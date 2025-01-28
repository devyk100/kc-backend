package types

import "sync"

type Jobs_t struct {
	List []Payload_t
	Mut  sync.RWMutex
}

type Payload_t struct {
	Key        string `json:"key"`
	Language   string `json:"lang"`
	QuestionId int    `json:"qid"`
	Token      string `json:"token"`
	Code       string `json:"code"`
}
