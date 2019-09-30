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
	{
		Version:     12,
		Description: "Alter roles with company_id",
		Script: `
ALTER TABLE roles
	ADD company_id  INT(10) UNSIGNED NOT NULL AFTER name,
	ADD KEY roles_company_id (company_id)
;`,
	},
	{
		Version:     13,
		Description: "Alter roles with companies constraint",
		Script: `
ALTER TABLE roles
	ADD CONSTRAINT fk_roles_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
;`,
	},
	{
		Version:     14,
		Description: "Add Products",
		Script: `
CREATE TABLE products (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	code	CHAR(10) NOT NULL UNIQUE,
	name	VARCHAR(255) NOT NULL,
	purchase_price DOUBLE NOT NULL DEFAULT 0,
	sale_price	DOUBLE NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY products_company_id (company_id),
	CONSTRAINT fk_products_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
	{
		Version:     15,
		Description: "Add Suppliers",
		Script: `
CREATE TABLE suppliers (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	code	CHAR(10) NOT NULL UNIQUE,
	name	VARCHAR(255) NOT NULL,
	address VARCHAR(255),
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY suppliers_company_id (company_id),
	CONSTRAINT fk_suppliers_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
	{
		Version:     16,
		Description: "Add Purchases",
		Script: `
CREATE TABLE purchases (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	supplier_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(10) NOT NULL UNIQUE,
	date	DATE NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY purchases_company_id (company_id),
	KEY purchases_branch_id (branch_id),
	KEY purchases_supplier_id (supplier_id),
	KEY purchases_created_by (created_by),
	KEY purchases_updated_by (updated_by),
	CONSTRAINT fk_purchases_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_purchases_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_purchases_to_suppliers FOREIGN KEY (supplier_id) REFERENCES suppliers(id),
	CONSTRAINT fk_purchases_to_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
	CONSTRAINT fk_purchases_to_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);`,
	},
	{
		Version:     17,
		Description: "Add Purchase Details",
		Script: `
CREATE TABLE purchase_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	purchase_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	price DOUBLE UNSIGNED NOT NULL,
	disc	DOUBLE UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY purchase_details_purchase_id (purchase_id),
	KEY purchase_details_product_id (product_id),
	CONSTRAINT fk_purchase_details_to_purchases FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_purchase_details_to_products FOREIGN KEY (product_id) REFERENCES products(id)
);`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
