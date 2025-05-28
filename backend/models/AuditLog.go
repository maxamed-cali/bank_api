package models

import "time"

type AuditLog struct {
	ID              uint
	UserID          *uint
	FullName        string
	AccountNumber   string
	ActionType      string
	TableName       string
	RecordID        uint
	Description     string
	ActionTimestamp time.Time
}
