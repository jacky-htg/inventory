package tests

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/jacky-htg/inventory/libraries/database/databasetest"
	"github.com/jacky-htg/inventory/schema"
)

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty.
//
// It does not return errors as this intended for testing only. Instead it will
// call Fatal on the provided testing.T if anything goes wrong.
//
// It returns the database to use as well as a function to call at the end of
// the test.
func NewUnit(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	databasetest.StartContainer(t)

	db, err := sql.Open(os.Getenv("DB_DRIVER"), "root:1234@tcp(localhost:33060)/rebel_db?parseTime=true")
	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}

	t.Log("waiting for database to be ready")

	// Wait for the database to be ready. Wait 100ms longer between each attempt.
	// Do not try more than 20 times.
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 1000 * time.Millisecond)
	}

	if pingError != nil {
		databasetest.StopContainer(t)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		db.Close()
		databasetest.StopContainer(t)
		t.Fatalf("migrating: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		databasetest.StopContainer(t)
	}

	return db, teardown
}
