package models

import "time"

type InfoUser struct {
	IsAdmin bool
	IsUser  bool
	Roles   []int
	Email   string
}

type WithCancelResponse struct {
	Response                  string    `json:"response"`
	StartTime                 time.Time `json:"handler_start_time"`
	EndTime                   time.Time `json:"handler_end_time"`
	ErrAfterTimeStartTime     time.Time `json:"err_after_time_start_time"`
	ErrAfterTimeEndTime       time.Time `json:"err_after_time_end_time"`
	ResponsefterTimeStartTime time.Time `json:"response_after_time_start_time"`
	ResponseAfterTimeEndTime  time.Time `json:"response_after_time_end_time"`
}

type TurnResponse struct {
	Turn int `json:"turn"`
}
