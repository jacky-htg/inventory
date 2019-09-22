package database

import (
	"context"
	"database/sql"
	"os"
)

//Open : open database
func Open() (*sql.DB, error) {
	return sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE"))
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sql.DB) error {
	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
