package services

import (
	"bank/db"
	"bank/models"
	"database/sql"
	"time"
)

// LogAudit inserts an audit log entry using raw SQL
func LogAudit(userID *uint, actionType, tableName string, recordID uint, description string) error {
	query := `
		INSERT INTO audit_logs (user_id, action_type, table_name, record_id, description, action_timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	// Use current timestamp
	timestamp := time.Now()

	// Handle nullable user_id
	var uid sql.NullInt64
	if userID != nil {
		uid = sql.NullInt64{Int64: int64(*userID), Valid: true}
	} else {
		uid = sql.NullInt64{Valid: false}
	}

	_, err := db.DB.Exec(query, uid, actionType, tableName, recordID, description, timestamp)
	return err
}

// GetAllAuditLogs returns all audit logs ordered by timestamp (DESC)
func GetAllAuditLogs() ([]models.AuditLog, error) {
	query := `
		SELECT id, user_id, action_type, table_name, record_id, description, action_timestamp
		FROM audit_logs
		ORDER BY action_timestamp DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog

	for rows.Next() {
		var log models.AuditLog
		var userID sql.NullInt64

		err := rows.Scan(
			&log.ID,
			&userID,
			&log.ActionType,
			&log.TableName,
			&log.RecordID,
			&log.Description,
			&log.ActionTimestamp,
		)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			uid := uint(userID.Int64)
			log.UserID = &uid
		} else {
			log.UserID = nil
		}

		logs = append(logs, log)
	}

	return logs, rows.Err()
}
