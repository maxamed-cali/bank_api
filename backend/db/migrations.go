package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func RunMigrations(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			full_name VARCHAR(100) NOT NULL,
			phone_number VARCHAR(15),
			address TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS credentials (
			id SERIAL PRIMARY KEY,
			user_id INTEGER UNIQUE NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_credential_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS user_roles (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			role_id INTEGER NOT NULL,
			CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
			CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id) ON UPDATE CASCADE ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			account_id VARCHAR(255) NOT NULL,
			transaction_type VARCHAR(20) NOT NULL,
			to_account_id VARCHAR(255),
			amount DECIMAL(15,2) NOT NULL,
			transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			description TEXT,
			CONSTRAINT fk_transaction_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			message TEXT NOT NULL,
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_notification_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS money_requests (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			requester_id VARCHAR(255) NOT NULL,
			recipient_id VARCHAR(255) NOT NULL,
			amount DECIMAL(15,2) NOT NULL,
			status VARCHAR(50),
			expires_at TIMESTAMP,
			requeste_at TIMESTAMP,
			recipient_user_id INTEGER NOT NULL,
			CONSTRAINT fk_money_request_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS audit_logs (
			id SERIAL PRIMARY KEY,
			user_id INTEGER,
			action_type VARCHAR(50) NOT NULL,
			table_name VARCHAR(100) NOT NULL,
			record_id INTEGER NOT NULL,
			action_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			description TEXT,
			CONSTRAINT fk_audit_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);`,

		`CREATE TABLE IF NOT EXISTS account_types (
			id SERIAL PRIMARY KEY,
			type_name VARCHAR(100) NOT NULL UNIQUE,
			description TEXT,
			currency VARCHAR(50) NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			account_number VARCHAR(255) NOT NULL UNIQUE,
			balance DECIMAL(15,2) DEFAULT 0.00,
			account_type_id INTEGER NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_account_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_account_type FOREIGN KEY (account_type_id) REFERENCES account_types(id) ON DELETE RESTRICT
		);`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	log.Println("All tables created successfully.")
	return nil
}
