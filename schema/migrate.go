package schema

import (
	"database/sql"

	"github.com/GuiaBolso/darwin"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
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
		Version:     2,
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
		Version:     3,
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
		Version:     4,
		Description: "Add Shelves",
		Script: `
CREATE TABLE shelves (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	branch_id	INT(10) UNSIGNED NOT NULL,
	code CHAR(10) NOT NULL,
	capacity MEDIUMINT(8) UNSIGNED NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY shelves_branch_id (branch_id),
	UNIQUE KEY shelves_code (branch_id, code),
	CONSTRAINT fk_shelves_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id)
);`,
	},
	{
		Version:     5,
		Description: "Add users",
		Script: `
CREATE TABLE users (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	username CHAR(15) NOT NULL UNIQUE,
	password VARCHAR(255) NOT NULL,
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
		Version:     6,
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
		Version:     7,
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
		Version:     8,
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
		Version:     9,
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
		Version:     10,
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
		Version:     11,
		Description: "Add Suppliers",
		Script: `
CREATE TABLE suppliers (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	code	CHAR(10) NOT NULL,
	name	VARCHAR(255) NOT NULL,
	address VARCHAR(255),
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY suppliers_company_id (company_id),
	UNIQUE KEY suppliers_code (company_id, code),
	CONSTRAINT fk_suppliers_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
	{
		Version:     12,
		Description: "Add Customers",
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
		Version:     13,
		Description: "Add Brands",
		Script: `
CREATE TABLE brands (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	code CHAR(10) NOT NULL,
	name VARCHAR(45) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	UNIQUE KEY brands_code (company_id, code),
	KEY brands_company_id (company_id),
	CONSTRAINT fk_brands_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
	{
		Version:     14,
		Description: "Add Categories",
		Script: `
CREATE TABLE categories (
	id   MEDIUMINT(8) UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(45) NOT NULL,
	PRIMARY KEY (id)
);`,
	},
	{
		Version:     15,
		Description: "Add Product Categories",
		Script: `
CREATE TABLE product_categories (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	category_id MEDIUMINT(8) UNSIGNED NOT NULL,
	name VARCHAR(45) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY product_categories_company_id (company_id),
	KEY product_categories_category_id (category_id),
	CONSTRAINT fk_product_categories_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_product_categories_to_categories FOREIGN KEY (category_id) REFERENCES categories(id)
);`,
	},
	{
		Version:     16,
		Description: "Add Products",
		Script: `
CREATE TABLE products (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	brand_id BIGINT(20) UNSIGNED NOT NULL,
	product_category_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(10) NOT NULL,
	name	VARCHAR(255) NOT NULL,
	purchase_price DOUBLE NOT NULL DEFAULT 0,
	sale_price	DOUBLE NOT NULL,
	minimum_stock MEDIUMINT(8) UNSIGNED NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY products_company_id (company_id),
	KEY products_brand_id (brand_id),
	KEY products_product_category_id (product_category_id),
	UNIQUE KEY products_code (company_id, code),
	CONSTRAINT fk_products_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_products_to_brands FOREIGN KEY (brand_id) REFERENCES brands(id),
	CONSTRAINT fk_products_to_product_categories FOREIGN KEY (product_category_id) REFERENCES product_categories(id)
);`,
	},
	{
		Version:     17,
		Description: "Add Purchases",
		Script: `
CREATE TABLE purchases (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	supplier_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL,
	date	DATE NOT NULL,
	disc DOUBLE NOT NULL DEFAULT 0,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY purchases_code (company_id, code),
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
		Version:     18,
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
		Version:     19,
		Description: "Add Purchase Returns",
		Script: `
CREATE TABLE purchase_returns (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	purchase_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL,
	date	DATE NOT NULL,
	disc DOUBLE NOT NULL DEFAULT 0,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY purchase_returns_code (company_id, code),
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
		Version:     20,
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
		Version:     21,
		Description: "Add Inventories",
		Script: `
CREATE TABLE inventories (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	product_code CHAR(20) NOT NULL,
	transaction_id BIGINT(20) UNSIGNED NOT NULL,
	code CHAR(20) NOT NULL,
	transaction_date DATE NOT NULL,
	type CHAR(2) NOT NULL,
	in_out TINYINT(1)UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL, 
	shelve_id BIGINT(20) UNSIGNED NOT NULL,
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
		Version:     22,
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
		Version:     23,
		Description: "Set Global log_bin_trust_function_creators",
		Script: `
SET GLOBAL log_bin_trust_function_creators = 1;
`,
	},
	{
		Version:     24,
		Description: "Add Good Receiving",
		Script: `
CREATE TABLE good_receivings (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	purchase_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL,
	date	DATE NOT NULL,
	remark VARCHAR(255) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY good_receivings_code (company_id, code),
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
		Version:     25,
		Description: "Add Good Receiving Details",
		Script: `
CREATE TABLE good_receiving_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	good_receiving_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	code CHAR(20) NOT NULL,
	shelve_id BIGINT(20) UNSIGNED NOT NULL,
	expired_date TIMESTAMP,
	PRIMARY KEY (id),
	KEY good_receiving_details_good_receiving_id (good_receiving_id),
	KEY good_receiving_details_product_id (product_id),
	KEY good_receiving_details_shelve_id (shelve_id),
	UNIQUE KEY good_receiving_details_code (code, product_id),
	CONSTRAINT fk_good_receiving_details_to_good_receivings FOREIGN KEY (good_receiving_id) REFERENCES good_receivings(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_good_receiving_details_to_products FOREIGN KEY (product_id) REFERENCES products(id),
	CONSTRAINT fk_good_receiving_details_to_shelves FOREIGN KEY (shelve_id) REFERENCES shelves(id)
);`,
	},
	{
		Version:     26,
		Description: "Add Saldo Stock Details",
		Script: `
CREATE TABLE saldo_stock_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	saldo_stock_id	BIGINT(20) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	code CHAR(20) NOT NULL,
	PRIMARY KEY (id),
	KEY saldo_stock_details_saldo_stock_id (saldo_stock_id),
	KEY saldo_stock_details_branch_id (branch_id),
	UNIQUE KEY saldo_stock_details_code (code, saldo_stock_id),
	CONSTRAINT fk_saldo_stock_details_to_saldo_stocks FOREIGN KEY (saldo_stock_id) REFERENCES saldo_stocks(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_saldo_stock_details_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id)
);`,
	},
	{
		Version:     27,
		Description: "Add Closing Stock Details",
		Script: `
CREATE PROCEDURE closing_stock_details(companyID int, curYear int, curMonth int)
BEGIN

	DECLARE nextYear int;
	DECLARE nextMonth int;
	
	IF curYear = 0 THEN
		SET curYear = year(now());
	END IF;
	
	IF curMonth = 0 THEN 
		SET curMonth = month(now());
	END IF;
	
	SET nextYear = curYear;
	SET nextMonth = curMonth + 1;
	
	IF curMonth = 12 THEN 
		SET nextYear = curYear+1;
		SET nextMonth = 1;
	END IF;
	
	INSERT INTO saldo_stock_details (saldo_stock_id, branch_id, code)
	SELECT 
		saldo_stocks.id saldo_stock_id,
		ifnull(inventories.branch_id, group_inventories.branch_id) branch_id,
		ifnull(inventories.product_code, group_inventories.product_code) code
	FROM (
		SELECT 
			MAX(union_inventories.id) id, 
			MAX(union_inventories.company_id) company_id, 
			MAX(union_inventories.branch_id) branch_id, 
			MAX(union_inventories.product_id) product_id, 
			MAX(union_inventories.product_code) product_code, 
			SUM(union_inventories.qty) qty
		FROM (
			(SELECT 
				0 id,
				saldo_stocks.company_id, 
				saldo_stock_details.branch_id,
				saldo_stocks.product_id, 
				saldo_stock_details.code product_code,
				1 qty  
			FROM saldo_stocks
			JOIN saldo_stock_details ON saldo_stocks.id = saldo_stock_details.saldo_stock_id
			WHERE saldo_stocks.year = curYear AND saldo_stocks.month = curMonth and saldo_stocks.company_id = companyID)
			union all
			(SELECT 
				inventories.id,
				inventories.company_id,
				inventories.branch_id,
				inventories.product_id,
				inventories.product_code,
				if(inventories.in_out, qty, -qty) as qty
			FROM inventories
			where month(inventories.transaction_date)=curMonth and year(inventories.transaction_date)=curYear and inventories.company_id = companyID)
		) union_inventories
		GROUP BY union_inventories.company_id, union_inventories.product_id, union_inventories.product_code
	) group_inventories
	left join inventories ON group_inventories.id = inventories.id and inventories.company_id = companyID
	join saldo_stocks ON ifnull(inventories.company_id, group_inventories.company_id)=saldo_stocks.company_id and ifnull(inventories.product_id, group_inventories.product_id)=saldo_stocks.product_id and saldo_stocks.year = nextYear and saldo_stocks.month = nextMonth and saldo_stocks.company_id = companyID
	WHERE group_inventories.qty > 0;

END;`,
	},
	{
		Version:     28,
		Description: "Add Closing Stock",
		Script: `
CREATE PROCEDURE closing_stocks(companyID int, curMonth int, curYear int)
BEGIN

	DECLARE nextYear int;
	DECLARE nextMonth int;
	
	IF curYear = 0 THEN
		SET curYear = year(now());
	END IF;
	
	IF curMonth = 0 THEN 
		SET curMonth = month(now());
	END IF;
	
	SET nextYear = curYear;
	SET nextMonth = curMonth + 1;
	
	IF curMonth = 12 THEN 
		SET nextYear = curYear+1;
		SET nextMonth = 1;
	END IF;
	
	INSERT INTO saldo_stocks (company_id, product_id, qty, year, month)
	SELECT companyID, saldo.id AS product_id, IF(transaction.qty IS NULL, saldo.qty, saldo.qty+transaction.qty) AS qty, nextYear AS year, nextMonth AS month 
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

call closing_stock_details(companyID, curYear, curMonth);

END`,
	},
	{
		Version:     29,
		Description: "Add Stock function",
		Script: `
CREATE FUNCTION stock (companyID int, productID int) RETURNS int(11)
BEGIN

declare stock int;

select SUM(union_stocks.qty) into stock
from (
	(select saldo_stocks.company_id, saldo_stocks.product_id, saldo_stocks.qty
	FROM saldo_stocks
	WHERE saldo_stocks.year=year(now()) and saldo_stocks.month=month(now()) and saldo_stocks.company_id = companyID and saldo_stocks.product_id = productID)
	union all
	(select inventories.company_id, inventories.product_id, if(inventories.in_out, qty, -qty) as qty 
	from inventories
	where month(inventories.transaction_date)=month(now()) and year(inventories.transaction_date)=year(now()) and inventories.company_id = companyID and product_id = productID )
) union_stocks
where union_stocks.company_id = companyID and union_stocks.product_id = productID
group by union_stocks.company_id, union_stocks.product_id;

RETURN stock;
END;`,
	},
	{
		Version:     30,
		Description: "Add Stock Branch function",
		Script: `
CREATE FUNCTION stock_branch(companyID int, branchID int, productID int) RETURNS int(11)
BEGIN

declare stock int;
declare curYear int;
declare curMonth int;

SET curYear = year(now());
SET curMonth = month(now());

select count(stocks.code) into stock
from ( 
	SELECT 
		ifnull(inventories.company_id, group_inventories.company_id) company_id,
		ifnull(inventories.product_id, group_inventories.product_id) product_id,
		ifnull(inventories.branch_id, group_inventories.branch_id) branch_id,
		ifnull(inventories.product_code, group_inventories.product_code) code
	FROM (
		SELECT 
			MAX(union_inventories.id) id, 
			MAX(union_inventories.company_id) company_id, 
			MAX(union_inventories.branch_id) branch_id, 
			MAX(union_inventories.product_id) product_id, 
			MAX(union_inventories.product_code) product_code, 
			SUM(union_inventories.qty) qty
		FROM (
			(SELECT 
				0 id,
				saldo_stocks.company_id, 
				saldo_stock_details.branch_id,
				saldo_stocks.product_id, 
				saldo_stock_details.code product_code,
				1 qty  
			FROM saldo_stocks
			JOIN saldo_stock_details ON saldo_stocks.id = saldo_stock_details.saldo_stock_id
			WHERE saldo_stocks.year = curYear AND saldo_stocks.month = curMonth and saldo_stocks.company_id = companyID)
			union all
			(SELECT 
				inventories.id,
				inventories.company_id,
				inventories.branch_id,
				inventories.product_id,
				inventories.product_code,
				if(inventories.in_out, qty, -qty) as qty
			FROM inventories
			where month(inventories.transaction_date)=curMonth and year(inventories.transaction_date)=curYear and inventories.company_id = companyID)
		) union_inventories
		GROUP BY union_inventories.company_id, union_inventories.product_id, union_inventories.product_code
	) group_inventories
	left join inventories ON group_inventories.id = inventories.id and inventories.company_id = companyID
		WHERE group_inventories.qty > 0
) stocks
where stocks.company_id = companyID and stocks.branch_id = branchID and stocks.product_id = productID
group by stocks.company_id, stocks.branch_id, stocks.product_id;

RETURN stock;
END;`,
	},
	{
		Version:     31,
		Description: "Add Stock Procedure",
		Script: `
CREATE PROCEDURE stocks(companyID int)
BEGIN

select union_stocks.product_id, SUM(union_stocks.qty) stock 
from (
	(select saldo_stocks.company_id, saldo_stocks.product_id, saldo_stocks.qty
	FROM saldo_stocks
	WHERE saldo_stocks.year=year(now()) and saldo_stocks.month=month(now()) and saldo_stocks.company_id=companyID)
	union all
	(select inventories.company_id, inventories.product_id, if(inventories.in_out, qty, -qty) as qty 
	from inventories
	where month(inventories.transaction_date)=month(now()) and year(inventories.transaction_date)=year(now()) and inventories.company_id=companyID )
) union_stocks
where union_stocks.company_id = companyID
group by union_stocks.company_id, union_stocks.product_id;

END;`,
	},
	{
		Version:     32,
		Description: "Add Branch Stocks Procedure",
		Script: `
CREATE PROCEDURE branch_stocks(companyID int, branchID int)
BEGIN

declare curYear int;
declare curMonth int;

SET curYear = year(now());
SET curMonth = month(now());

select stocks.product_id, count(stocks.code) stock
from ( 
	SELECT 
		ifnull(inventories.company_id, group_inventories.company_id) company_id,
		ifnull(inventories.product_id, group_inventories.product_id) product_id,
		ifnull(inventories.branch_id, group_inventories.branch_id) branch_id,
		ifnull(inventories.product_code, group_inventories.product_code) code
	FROM (
		SELECT 
			MAX(union_inventories.id) id, 
			MAX(union_inventories.company_id) company_id, 
			MAX(union_inventories.branch_id) branch_id, 
			MAX(union_inventories.product_id) product_id, 
			MAX(union_inventories.product_code) product_code, 
			SUM(union_inventories.qty) qty
		FROM (
			(SELECT 
				0 id,
				saldo_stocks.company_id, 
				saldo_stock_details.branch_id,
				saldo_stocks.product_id, 
				saldo_stock_details.code product_code,
				1 qty  
			FROM saldo_stocks
			JOIN saldo_stock_details ON saldo_stocks.id = saldo_stock_details.saldo_stock_id
			WHERE saldo_stocks.year = curYear AND saldo_stocks.month = curMonth and saldo_stocks.company_id = companyID)
			union all
			(SELECT 
				inventories.id,
				inventories.company_id,
				inventories.branch_id,
				inventories.product_id,
				inventories.product_code,
				if(inventories.in_out, qty, -qty) as qty
			FROM inventories
			where month(inventories.transaction_date)=curMonth and year(inventories.transaction_date)=curYear and inventories.company_id = companyID)
		) union_inventories
		GROUP BY union_inventories.company_id, union_inventories.product_id, union_inventories.product_code
	) group_inventories
	left join inventories ON group_inventories.id = inventories.id and inventories.company_id = companyID
		WHERE group_inventories.qty > 0
) stocks
where stocks.company_id = companyID and stocks.branch_id = branchID
group by stocks.company_id, stocks.branch_id, stocks.product_id;

END;`,
	},
	{
		Version:     33,
		Description: "Add Branch Stock details Procedure",
		Script: `
CREATE PROCEDURE branch_stock_details(companyID int, branchID int, productID int)
BEGIN

declare curYear int;
declare curMonth int;

SET curYear = year(now());
SET curMonth = month(now());

select stocks.product_id, stocks.code
from ( 
	SELECT 
		ifnull(inventories.company_id, group_inventories.company_id) company_id,
		ifnull(inventories.product_id, group_inventories.product_id) product_id,
		ifnull(inventories.branch_id, group_inventories.branch_id) branch_id,
		ifnull(inventories.product_code, group_inventories.product_code) code
	FROM (
		SELECT 
			MAX(union_inventories.id) id, 
			MAX(union_inventories.company_id) company_id, 
			MAX(union_inventories.branch_id) branch_id, 
			MAX(union_inventories.product_id) product_id, 
			MAX(union_inventories.product_code) product_code, 
			SUM(union_inventories.qty) qty
		FROM (
			(SELECT 
				0 id,
				saldo_stocks.company_id, 
				saldo_stock_details.branch_id,
				saldo_stocks.product_id, 
				saldo_stock_details.code product_code,
				1 qty  
			FROM saldo_stocks
			JOIN saldo_stock_details ON saldo_stocks.id = saldo_stock_details.saldo_stock_id
			WHERE saldo_stocks.year = curYear AND saldo_stocks.month = curMonth and saldo_stocks.company_id = companyID)
			union all
			(SELECT 
				inventories.id,
				inventories.company_id,
				inventories.branch_id,
				inventories.product_id,
				inventories.product_code,
				if(inventories.in_out, qty, -qty) as qty
			FROM inventories
			where month(inventories.transaction_date)=curMonth and year(inventories.transaction_date)=curYear and inventories.company_id = companyID)
		) union_inventories
		GROUP BY union_inventories.company_id, union_inventories.product_id, union_inventories.product_code
	) group_inventories
	left join inventories ON group_inventories.id = inventories.id and inventories.company_id = companyID
		WHERE group_inventories.qty > 0
) stocks
where stocks.company_id = companyID and stocks.branch_id = branchID;

END;`,
	},
	{
		Version:     34,
		Description: "Add Return Receive",
		Script: `
CREATE TABLE receiving_returns (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(20) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	good_receiving_id BIGINT(20) UNSIGNED NOT NULL,
	date DATE NOT NULL,
	code CHAR(13) NOT NULL,
	remark VARCHAR(255) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL, 
	PRIMARY KEY (id),
	KEY receiving_returns_company_id (company_id),
	KEY receiving_returns_branch_id (branch_id),
	KEY receiving_returns_good_receiving_id (good_receiving_id),
	KEY receiving_returns_created_by (created_by),
	KEY receiving_returns_updated_by (updated_by),
	UNIQUE KEY receiving_returns_code (code, company_id),
	CONSTRAINT fk_receiving_returns_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_receiving_returns_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_receiving_returns_to_good_receivings FOREIGN KEY (good_receiving_id) REFERENCES good_receivings(id),
	CONSTRAINT fk_receiving_returns_to_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
	CONSTRAINT fk_receiving_returns_to_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);`,
	},
	{
		Version:     35,
		Description: "Add Receiving Return Details",
		Script: `
CREATE TABLE receiving_return_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	receiving_return_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	code CHAR(20) NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY receiving_return_details_receiving_return_id (receiving_return_id),
	KEY receiving_return_details_product_id (product_id),
	UNIQUE KEY receiving_return_details_code (code, product_id),
	CONSTRAINT fk_receiving_return_details_to_receiving_returns FOREIGN KEY (receiving_return_id) REFERENCES receiving_returns(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_receiving_return_details_to_products FOREIGN KEY (product_id) REFERENCES products(id)
);`,
	},
	{
		Version:     36,
		Description: "Add Salesmen",
		Script: `
CREATE TABLE salesmen (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	code	CHAR(10) NOT NULL UNIQUE,
	name VARCHAR(100) NOT NULL,
	email VARCHAR(100) NOT NULL UNIQUE,
	address VARCHAR(255) NOT NULL,
	hp CHAR(15) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	KEY salesmen_company_id (company_id),
	CONSTRAINT fk_salesmen_to_companies FOREIGN KEY (company_id) REFERENCES companies(id)
);`,
	},
	{
		Version:     37,
		Description: "Add SalesOrders",
		Script: `
CREATE TABLE sales_orders (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	salesman_id BIGINT(20) UNSIGNED NOT NULL,
	customer_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL,
	date	DATE NOT NULL,
	disc DOUBLE NOT NULL DEFAULT 0,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY sales_orders_code (company_id, code),
	KEY sales_orders_company_id (company_id),
	KEY sales_orders_branch_id (branch_id),
	KEY sales_orders_salesman_id (salesman_id),
	KEY sales_orders_customer_id (customer_id),
	KEY sales_orders_created_by (created_by),
	KEY sales_orders_updated_by (updated_by),
	CONSTRAINT fk_sales_orders_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_sales_orders_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_sales_orders_to_salesmen FOREIGN KEY (salesman_id) REFERENCES salesmen(id),
	CONSTRAINT fk_sales_orders_to_customers FOREIGN KEY (customer_id) REFERENCES customers(id),
	CONSTRAINT fk_sales_orders_to_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
	CONSTRAINT fk_sales_orders_to_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);`,
	},
	{
		Version:     38,
		Description: "Add SalesOrders Details",
		Script: `
CREATE TABLE sales_order_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	sales_order_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	price DOUBLE UNSIGNED NOT NULL,
	disc	DOUBLE UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY sales_order_details_sales_order_id (sales_order_id),
	KEY sales_order_details_product_id (product_id),
	CONSTRAINT fk_sales_order_details_to_sales_orders FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_sales_order_details_to_products FOREIGN KEY (product_id) REFERENCES products(id)
);`,
	},
	{
		Version:     39,
		Description: "Add SalesOrder Returns",
		Script: `
CREATE TABLE sales_order_returns (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	sales_order_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL,
	date	DATE NOT NULL,
	disc DOUBLE NOT NULL DEFAULT 0,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY sales_order_returns_code (company_id, code),
	KEY sales_order_returns_company_id (company_id),
	KEY sales_order_returns_branch_id (branch_id),
	KEY sales_order_returns_sales_order_id (sales_order_id),
	KEY sales_order_returns_created_by (created_by),
	KEY sales_order_returns_updated_by (updated_by),
	CONSTRAINT fk_sales_order_returns_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_sales_order_returns_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_sales_order_returns_to_sales_orders FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id),
	CONSTRAINT fk_sales_order_returns_to_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
	CONSTRAINT fk_sales_order_returns_to_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);`,
	},
	{
		Version:     40,
		Description: "Add SalesOrder Return Details",
		Script: `
CREATE TABLE sales_order_return_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	sales_order_return_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	price DOUBLE UNSIGNED NOT NULL,
	disc	DOUBLE UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY sales_order_return_details_sales_order_return_id (sales_order_return_id),
	KEY sales_order_return_details_product_id (product_id),
	CONSTRAINT fk_sales_order_return_details_to_sales_order_returns FOREIGN KEY (sales_order_return_id) REFERENCES sales_order_returns(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_sales_order_return_details_to_products FOREIGN KEY (product_id) REFERENCES products(id)
);`,
	},
	{
		Version:     41,
		Description: "Add Delivery Order",
		Script: `
CREATE TABLE deliveries (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(10) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	sales_order_id BIGINT(20) UNSIGNED NOT NULL,
	code	CHAR(13) NOT NULL,
	date	DATE NOT NULL,
	remark VARCHAR(255) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY deliveries_code (company_id, code),
	KEY deliveries_company_id (company_id),
	KEY deliveries_branch_id (branch_id),
	KEY deliveries_sales_order_id (sales_order_id),
	KEY deliveries_created_by (created_by),
	KEY deliveries_updated_by (updated_by),
	CONSTRAINT fk_deliveries_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_deliveries_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_deliveries_to_sales_orders FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id),
	CONSTRAINT fk_deliveries_to_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
	CONSTRAINT fk_deliveries_to_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);`,
	},
	{
		Version:     42,
		Description: "Add Delivery Details",
		Script: `
CREATE TABLE delivery_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	delivery_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	code CHAR(20) NOT NULL,
	shelve_id BIGINT(20) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY delivery_details_delivery_id (delivery_id),
	KEY delivery_details_product_id (product_id),
	KEY delivery_details_shelve_id (shelve_id),
	UNIQUE KEY delivery_details_code (code, product_id),
	CONSTRAINT fk_delivery_details_to_deliveries FOREIGN KEY (delivery_id) REFERENCES deliveries(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_delivery_details_to_products FOREIGN KEY (product_id) REFERENCES products(id),
	CONSTRAINT fk_delivery_details_to_shelves FOREIGN KEY (shelve_id) REFERENCES shelves(id)
);`,
	},
	{
		Version:     43,
		Description: "Add Return Delivery",
		Script: `
CREATE TABLE delivery_returns (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	company_id	INT(20) UNSIGNED NOT NULL,
	branch_id INT(10) UNSIGNED NOT NULL,
	delivery_id BIGINT(20) UNSIGNED NOT NULL,
	date DATE NOT NULL,
	code CHAR(13) NOT NULL,
	remark VARCHAR(255) NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	updated TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by BIGINT(20) UNSIGNED NOT NULL,
	updated_by BIGINT(20) UNSIGNED NOT NULL, 
	PRIMARY KEY (id),
	KEY delivery_returns_company_id (company_id),
	KEY delivery_returns_branch_id (branch_id),
	KEY delivery_returns_delivery_id (delivery_id),
	KEY delivery_returns_created_by (created_by),
	KEY delivery_returns_updated_by (updated_by),
	UNIQUE KEY delivery_returns_code (code, company_id),
	CONSTRAINT fk_delivery_returns_to_companies FOREIGN KEY (company_id) REFERENCES companies(id),
	CONSTRAINT fk_delivery_returns_to_branches FOREIGN KEY (branch_id) REFERENCES branches(id),
	CONSTRAINT fk_delivery_returns_to_deliveries FOREIGN KEY (delivery_id) REFERENCES deliveries(id),
	CONSTRAINT fk_delivery_returns_to_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
	CONSTRAINT fk_delivery_returns_to_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);`,
	},
	{
		Version:     44,
		Description: "Add Delivery Return Details",
		Script: `
CREATE TABLE delivery_return_details (
	id   BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	delivery_return_id	BIGINT(20) UNSIGNED NOT NULL,
	product_id BIGINT(20) UNSIGNED NOT NULL,
	code CHAR(20) NOT NULL,
	qty MEDIUMINT(8) UNSIGNED NOT NULL,
	PRIMARY KEY (id),
	KEY delivery_return_details_delivery_return_id (delivery_return_id),
	KEY delivery_return_details_product_id (product_id),
	UNIQUE KEY delivery_return_details_code (code, product_id),
	CONSTRAINT fk_delivery_return_details_to_delivery_returns FOREIGN KEY (delivery_return_id) REFERENCES delivery_returns(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_delivery_return_details_to_products FOREIGN KEY (product_id) REFERENCES products(id)
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
