package handlers

import "github.com/maeshinshin/go-multiapi/internal/database"

type Handlers struct {
	db database.Service
}

func NewHandlers(db database.Service) *Handlers {
	return &Handlers{db}
}

type HelloWorldresponse struct {
	Message      string  `json:"message"`
	Message2     string  `json:"message2"`
	Message3     string  `json:"message3"`
	ErrorMessage string  `json:"errormessage,omitempty"`
	Error        *Errors `json:"errors"`
}

type Errors struct {
	Error  string `json:"error"`
	Error2 string `json:"error2"`
}

type HealthResponse struct {
	Status             string `json:"status"`
	Message            string `json:"message"`
	Openconnections    string `json:"open_connections"`
	Inuse              string `json:"in_use"`
	Idle               string `json:"idle"`
	Waitcount          string `json:"wait_count"`
	Waitduration       string `json:"wait_duration"`
	Maxidle_closed     string `json:"max_idle_closed"`
	Maxlifetime_closed string `json:"max_lifetime_closed"`
	Error              string `json:"error,omitempty"`
}
