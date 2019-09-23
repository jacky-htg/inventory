package schema

import (
	"database/sql"

	"github.com/GuiaBolso/darwin"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add users",
		Script: `
CREATE TABLE users (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	username         CHAR(15) NOT NULL UNIQUE,
	password         varchar(255) NOT NULL,
	email     VARCHAR(255) NOT NULL UNIQUE,
	is_active TINYINT(1) NOT NULL DEFAULT '0',
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id)
);`,
	},
	{
		Version:     2,
		Description: "Add access",
		Script: `
CREATE TABLE access (
	id   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	parent_id         INT(10) UNSIGNED,
	name         varchar(255) NOT NULL UNIQUE,
	alias         varchar(255) NOT NULL UNIQUE,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id)
);`,
	},
	{
		Version:     3,
		Description: "Add roles",
		Script: `
CREATE TABLE roles (
	id   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	name         varchar(255) NOT NULL UNIQUE,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id)
);`,
	},
	{
		Version:     4,
		Description: "Add access_roles",
		Script: `
CREATE TABLE access_roles (
	id   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	access_id         INT(10) UNSIGNED NOT NULL,
	role_id         INT(10) UNSIGNED NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	UNIQUE KEY access_roles_unique (access_id, role_id),
	KEY access_roles_access_id (access_id),
	KEY access_roles_role_id (role_id),
	CONSTRAINT fk_access_roles_to_access FOREIGN KEY (access_id) REFERENCES access(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_access_roles_to_roles FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE
);`,
	},
	{
		Version:     5,
		Description: "Add roles_users",
		Script: `
CREATE TABLE roles_users (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	role_id         INT(10) UNSIGNED NOT NULL,
	user_id         BIGINT(20) UNSIGNED NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	UNIQUE KEY roles_users_unique (role_id, user_id),
	KEY roles_users_role_id (role_id),
	KEY roles_users_user_id (user_id),
	CONSTRAINT fk_roles_users_to_roles FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_roles_users_to_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);`,
	},
	{
		Version:     6,
		Description: "Add companies",
		Script: `
CREATE TABLE companies (
	id   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	code	char(10) NOT NULL UNIQUE,
	name	varchar(255) NOT NULL,
	address	varchar(255),
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id)
);`,
	},
	{
		Version:     7,
		Description: "Add regions",
		Script: `
CREATE TABLE regions (
	id   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id   INT(10) UNSIGNED NOT NULL,
	code	char(10) NOT NULL UNIQUE,
	name	varchar(255) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY regions_company_id (company_id),
	CONSTRAINT fk_regions_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
	{
		Version:     8,
		Description: "Add branches",
		Script: `
CREATE TABLE branches (
	id   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	code	char(10) NOT NULL UNIQUE,
	name	varchar(255) NOT NULL,
	type char(1) NOT NULL,
	address	varchar(255),
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY branches_company_id (company_id),
	CONSTRAINT fk_branches_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
	{
		Version:     9,
		Description: "Add branches_regions",
		Script: `
CREATE TABLE branches_regions (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	branch_id	INT(10) UNSIGNED NOT NULL,
	region_id	INT(10) UNSIGNED NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	UNIQUE KEY branches_regions_unique (branch_id, region_id),
	KEY branches_regions_branch_id (branch_id),
	KEY branches_regions_region_id (region_id),
	CONSTRAINT fk_branches_regions_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_branches_regions_to_regions FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE ON UPDATE CASCADE
);`,
	},
	{
		Version:     10,
		Description: "Alter users with company_id, region_id and branch_id",
		Script: `
ALTER TABLE users
	ADD company_id  INT(10) UNSIGNED NOT NULL AFTER is_active,
	ADD region_id INT(10) UNSIGNED AFTER company_id,
	ADD branch_id INT(10) UNSIGNED AFTER region_id,
	ADD KEY users_company_id (company_id),
	ADD KEY users_region_id (region_id),
	ADD KEY users_branch_id (branch_id),
	ADD CONSTRAINT fk_users_to_regions FOREIGN KEY (region_id) REFERENCES regions(id),
	ADD CONSTRAINT fk_users_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id)
;`,
	},
	{
		Version:     11,
		Description: "Alter users with companies constraint",
		Script: `
ALTER TABLE users
	ADD CONSTRAINT fk_users_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
;`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
