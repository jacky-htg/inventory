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
	company_id  INT(10) UNSIGNED NOT NULL,
	region_id INT(10) UNSIGNED,
	branch_id INT(10) UNSIGNED,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY users_company_id (company_id),
	KEY users_region_id (region_id),
	KEY users_branch_id (branch_id),
	CONSTRAINT fk_users_to_regions FOREIGN KEY (region_id) REFERENCES regions(id),
	CONSTRAINT fk_users_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_users_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
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
	company_id  INT(10) UNSIGNED NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY roles_company_id (company_id),
	CONSTRAINT fk_roles_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
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
		Version:     11,
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
		Version:     12,
		Description: "Add Purchases",
		Script: `
CREATE TABLE purchases (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	supplier_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL UNIQUE,
	date	DATE NOT NULL,
	disc DOUBLE NOT NULL DEFAULT 0,
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
		Version:     13,
		Description: "Add Purchase Details",
		Script: `
CREATE TABLE purchase_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	purchase_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	price DOUBLE UNSIGNED NOT NULL,
	disc	DOUBLE UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY purchase_details_purchase_id (purchase_id),
	KEY purchase_details_product_id (product_id),
	CONSTRAINT fk_purchase_details_to_purchases FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_purchase_details_to_products FOREIGN KEY (product_id) REFERENCES products(id)
);`,
	},
	{
		Version:     14,
		Description: "Add Purchase Returns",
		Script: `
CREATE TABLE purchase_returns (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	purchase_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL UNIQUE,
	date	DATE NOT NULL,
	disc DOUBLE NOT NULL DEFAULT 0,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY purchase_returns_company_id (company_id),
	KEY purchase_returns_branch_id (branch_id),
	KEY purchase_returns_purchase_id (purchase_id),
	KEY purchase_returns_created_by (created_by),
	KEY purchase_returns_updated_by (updated_by),
	CONSTRAINT fk_purchase_returns_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_purchase_returns_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_purchase_returns_to_purchases FOREIGN KEY (purchase_id) REFERENCES purchases(id),
	CONSTRAINT fk_purchase_returns_to_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
	CONSTRAINT fk_purchase_returns_to_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);`,
	},
	{
		Version:     15,
		Description: "Add Purchase Return Details",
		Script: `
CREATE TABLE purchase_return_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	purchase_return_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	price DOUBLE UNSIGNED NOT NULL,
	disc	DOUBLE UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY purchase_return_details_purchase_return_id (purchase_return_id),
	KEY purchase_return_details_product_id (product_id),
	CONSTRAINT fk_purchase_return_details_to_purchase_returns FOREIGN KEY (purchase_return_id) REFERENCES purchase_returns(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_purchase_return_details_to_products FOREIGN KEY (product_id) REFERENCES products(id)
);`,
	},
	{
		Version:     16,
		Description: "Add Inventories",
		Script: `
CREATE TABLE inventories (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	transaction_id BIGINT(20) UNSIGNED NOT NULL,
	transaction_date DATE NOT NULL,
	type CHAR(2) NOT NULL,
	in_out TINYINT(1)UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL, 
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY inventories_company_id (company_id),
	KEY inventories_product_id (product_id),
	CONSTRAINT fk_inventories_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_inventories_to_products FOREIGN KEY (product_id) REFERENCES products(id)
);`,
	},
	{
		Version:     17,
		Description: "Add saldo stocks",
		Script: `
CREATE TABLE saldo_stocks (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	year YEAR(4) NOT NULL,
	month TINYINT(2) NOT NULL,
	PRIMARY KEY (id),
	KEY saldo_stocks_company_id (company_id),
	KEY saldo_stocks_product_id (product_id),
	CONSTRAINT fk_saldo_stocks_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_saldo_stocks_to_products FOREIGN KEY (product_id) REFERENCES products(id)
);`,
	},
	{
		Version:     18,
		Description: "Set Global log_bin_trust_function_creators",
		Script: `
SET GLOBAL log_bin_trust_function_creators = 1;
`,
	},
	{
		Version:     19,
		Description: "Add Closing Stocks",
		Script: `
CREATE FUNCTION closing_stocks(companyID int, curMonth int, curYear int)
RETURNS INTEGER
BEGIN

	DECLARE nextYear int;
	DECLARE nextMonth int;
	
	SET nextYear = curYear;
	SET nextMonth = curMonth + 1;
	
	IF curMonth = 12 THEN 
		SET nextYear = curYear+1;
		SET nextMonth = 1;
	END IF;
	
	INSERT INTO saldo_stocks (company_id, product_id, qty, year, month)
	SELECT saldo.id AS product_id, IF(transaction.qty IS NULL, saldo.qty, saldo.qty+transaction.qty) AS qty, nextYear AS year, nextMonth AS month 
	FROM (
		SELECT products.id, IF(saldo_stocks.qty IS NULL, 0, saldo_stocks.qty) AS qty
		FROM products
		LEFT JOIN saldo_stocks ON products.id = saldo_stocks.product_id AND products.company_id=saldo_stocks.company_id AND saldo_stocks.year=curYear AND saldo_stocks.month=curMonth
		WHERE products.company_id=companyID
	) AS saldo
	LEFT JOIN (
		SELECT tr.product_id, SUM(tr.qty) AS qty
		FROM (
			select inventories.product_id, if(inventories.in_out, qty, -qty) as qty   
			from inventories
			WHERE MONTH(inventories.transaction_date)=curMonth AND YEAR(inventories.transaction_date)=curYear AND inventories.company_id=companyID
		) as tr
		GROUP BY tr.product_id
	) as transaction ON saldo.id=transaction.product_id;

RETURN 1;
END;`,
	},
	{
		Version:     20,
		Description: "Add Closing Stocks",
		Script: `
CREATE TABLE customers (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	name VARCHAR(100) NOT NULL,
	email VARCHAR(100) NOT NULL UNIQUE,
	address VARCHAR(255) NOT NULL,
	hp CHAR(15) NOT NULL,
	PRIMARY KEY (id),
	KEY customers_company_id (company_id),
	CONSTRAINT fk_customers_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
	{
		Version:     21,
		Description: "Add Good Receiving",
		Script: `
CREATE TABLE good_receivings (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	purchase_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL UNIQUE,
	date	DATE NOT NULL,
	remark VARCHAR(255) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY good_receivings_company_id (company_id),
	KEY good_receivings_branch_id (branch_id),
	KEY good_receivings_purchase_id (purchase_id),
	KEY good_receivings_created_by (created_by),
	KEY good_receivings_updated_by (updated_by),
	CONSTRAINT fk_good_receivings_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_good_receivings_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_good_receivings_to_purchases FOREIGN KEY (purchase_id) REFERENCES purchases(id),
	CONSTRAINT fk_good_receivings_to_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
	CONSTRAINT fk_good_receivings_to_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);`,
	},
	{
		Version:     22,
		Description: "Add Good Receiving Details",
		Script: `
CREATE TABLE good_receiving_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	good_receiving_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY good_receiving_details_good_receiving_id (good_receiving_id),
	KEY good_receiving_details_product_id (product_id),
	CONSTRAINT fk_good_receiving_details_to_good_receivings FOREIGN KEY (good_receiving_id) REFERENCES good_receivings(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_good_receiving_details_to_products FOREIGN KEY (product_id) REFERENCES products(id)
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
