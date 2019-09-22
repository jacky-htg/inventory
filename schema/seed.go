package schema

import (
	"database/sql"
	"fmt"
)

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Using a constant in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.

const seedUsers string = `
INSERT INTO users (id, username, password, email, is_active) VALUES
	(1, 'jackyhtg', '$2y$10$ekouPwVdtMEy5AFbogzfSeRloxHzUwEAsM7SyNJXnso/F9ds/XUYy', 'admin@admin.com', 1);
`

const seedAccess string = `
INSERT INTO access (id, name, alias, created) VALUES (1, 'root', 'root', NOW());
`

const seedRoles string = `
INSERT INTO roles (id, name, created) VALUES (1, 'superadmin', NOW());
`

const seedAccessRoles string = `
INSERT INTO access_roles (access_id, role_id) VALUES (1, 1);
`

const seedRolesUsers string = `
INSERT INTO roles_users (role_id, user_id) VALUES (1, 1);
`

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sql.DB) error {
	seeds := []string{
		seedUsers,
		seedAccess,
		seedRoles,
		seedAccessRoles,
		seedRolesUsers,
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, seed := range seeds {
		_, err = tx.Exec(seed)
		if err != nil {
			tx.Rollback()
			fmt.Println("error execute seed")
			return err
		}
	}

	return tx.Commit()
}
