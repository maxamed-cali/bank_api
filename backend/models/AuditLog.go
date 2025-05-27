package models

import "time"

type AuditLog struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    UserID          *uint     `json:"user_id"`
    ActionType      string    `json:"action_type"`  // CREATE, UPDATE, DELETE
    TableName       string    `json:"table_name"`
    RecordID        uint      `json:"record_id"`
    ActionTimestamp time.Time `gorm:"autoCreateTime" json:"action_timestamp"`
    Description     string    `json:"description"`
}
