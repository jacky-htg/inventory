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
INSERT INTO users (id, username, password, email, is_active, company_id, branch_id) VALUES 
(1, 'jackyhtg', '$2y$10$ekouPwVdtMEy5AFbogzfSeRloxHzUwEAsM7SyNJXnso/F9ds/XUYy', 'admin@admin.com', 1, 1, null),
(2, 'peterpan', '$2a$10$gT5pAqbiLxXTElwluvNJuef0jlOlHyt4q9ApC7jyhMb49OvIHeKgO', 'peterpan@gmail.com', 1, 1, 1);
`

const seedAccess string = `
INSERT INTO access (id, name, alias, created) VALUES (1, 'root', 'root', NOW());
`

const seedRoles string = `
INSERT INTO roles (id, name, company_id, created) VALUES (1, 'superadmin', 1, NOW());
`

const seedAccessRoles string = `
INSERT INTO access_roles (access_id, role_id) VALUES (1, 1);
`

const seedRolesUsers string = `
INSERT INTO roles_users (role_id, user_id) VALUES 
(1, 1),
(1, 2);
`

const seedCompanies string = `
INSERT INTO companies (id, code, name) VALUES (1, "DM", "Dummy");
`

const seedCategories string = `
INSERT INTO categories (id, name) VALUES (1, "Accesories");
`

const seedBranches string = `
INSERT INTO branches (id, company_id, code, name, type, address ) VALUES (1, 1, "123", "Toko Bagus", "s", "jalan jalan");
`
const seedShelves string = `
INSERT INTO shelves (id, branch_id, code, capacity) VALUES (1, 1, "SHV-01", 1000);
`

const seedSuppliers string = `
INSERT INTO suppliers (id, company_id, code, name, address ) VALUES (1, 1, "SUP_01", "Supplier Test", "jalan supplier");
`

const seedBrands string = `
INSERT INTO brands (id, company_id, code, name) VALUES (1, 1, "BRAND-01", "Brand Test");
`

const seedProductCategories string = `
INSERT INTO product_categories (id, company_id, category_id, name) VALUES(1, 1, 1, "Furniture");
`

const seedProducts string = `
INSERT INTO products (id, company_id, brand_id, product_category_id, code, name, sale_price, minimum_stock) VALUES 
(1, 1, 1, 1, "PROD-01", "Product Satu", "1000", "25"),
(2, 1, 1, 1, "PROD-02", "Product Dua", "500", "1000");
`

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sql.DB) error {
	seeds := []string{
		seedCompanies,
		seedBranches,
		seedUsers,
		seedAccess,
		seedRoles,
		seedAccessRoles,
		seedRolesUsers,
		seedCategories,
		seedShelves,
		seedSuppliers,
		seedBrands,
		seedProductCategories,
		seedProducts,
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
