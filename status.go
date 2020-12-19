package zrule

import "time"

type ServerStatus struct {
	Players       uint      `json:"players"`
	ServerVersion string    `json:"server_version"`
	StartTime     time.Time `json:"start_time"`
}
