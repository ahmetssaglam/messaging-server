package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"messaging-server/internal/configs"
	log "messaging-server/internal/logging"
	"messaging-server/internal/models"
)

type PostgresDB struct {
	*sql.DB
}

var PostgresConnection *PostgresDB

const fetchQuery = `
        SELECT id, content, phone_number, is_sent
        FROM messages
        WHERE is_sent = FALSE
        ORDER BY id	
        LIMIT $1
    `

const fetchAllSentQuery = `
    SELECT id, content, phone_number, is_sent
  	FROM messages
    WHERE is_sent = TRUE
  	ORDER BY id
`

const updateQuery = `
        UPDATE messages
           SET is_sent = TRUE
         WHERE id = $1
    `

// ConnectPostgres initializes DB on first call; returns an error if it fails.
// Subsequent calls are no-ops.
func ConnectPostgres() error {
	if PostgresConnection != nil {
		return nil
	}

	db, err := sql.Open("postgres", configs.PostgresConfig.ConnStr)
	if err != nil {
		return fmt.Errorf("failed to open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	PostgresConnection = &PostgresDB{DB: db}
	log.Logger.Info("Postgres connection established")
	return nil
}

// ensureConnection pings the DB and, if it has dropped, tears it down
// and re-runs ConnectPostgres to rebuild the singleton.
func (p *PostgresDB) ensureConnection() {
	if err := p.Ping(); err != nil {
		log.Logger.Warningf("lost DB connection (%v), reconnecting…", err)
		_ = p.Close()

		// clear DB so ConnectPostgres will re-init
		PostgresConnection = nil
		if err := ConnectPostgres(); err != nil {
			log.Logger.Fatalf("reconnect failed: %v", err)
		}
		// reset receiver to the fresh global
		*p = *PostgresConnection
	}
}

func (p *PostgresDB) FetchPendingMessages(limit int) ([]models.Message, error) {

	p.ensureConnection()

	// execute the query
	rows, err := p.Query(fetchQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("query pending messages: %w", err)
	}
	defer rows.Close()

	var msgs []models.Message

	// iterate over the rows
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.Content, &m.PhoneNumber, &m.IsSent); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	log.Logger.Debug("Messages fetched from DB")
	return msgs, nil
}

func (p *PostgresDB) MarkSent(id string) error {

	p.ensureConnection()

	// execute the update query
	res, err := p.Exec(updateQuery, id)
	if err != nil {
		return fmt.Errorf("updating is_sent for id %s: %w", id, err)
	}

	// check if any rows were affected
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no message found with id %s", id)
	}

	log.Logger.Debugf("Message %s marked as sent", id)
	return nil
}

// FetchAllSentMessages retrieves *all* messages that have been sent.
func (p *PostgresDB) FetchAllSentMessages() ([]models.Message, error) {
	p.ensureConnection()

	rows, err := p.Query(fetchAllSentQuery)
	if err != nil {
		return nil, fmt.Errorf("query all pending messages: %w", err)
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.Content, &m.PhoneNumber, &m.IsSent); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	log.Logger.Debugf("Fetched %d unsent messages", len(msgs))
	return msgs, nil
}
