package models

import "database/sql"

//Branch : struct of Branch
type Branch struct {
	ID      uint32
	Code    string
	Name    string
	Address sql.NullString
	Type    string
	Company Company
}
