package models

import "time"

type OutboxEvent struct {
    ID        int64      `json:"id"`
    Channel   string     `json:"channel"`
    Payload   string     `json:"payload"`
    Attempts  int        `json:"attempts"`
    LastError string     `json:"last_error,omitempty"`
    CreatedAt time.Time  `json:"created_at"`
    SentAt    *time.Time `json:"sent_at,omitempty"`
}
