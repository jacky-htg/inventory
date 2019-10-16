package api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// GetCode is function to generate code of transaction
func GetCode(ctx context.Context, tx *sql.Tx, prefix string, tableName string, companyID uint32) (string, error) {
	var code string
	var codeInt int

	if len(prefix) > 0 {
		runes := []rune(prefix)
		prefix = string(runes[0:2])
	}

	prefix += time.Now().Format("200601")

	query := "SELECT code FROM " + tableName + " WHERE company_id = ? AND code LIKE ? ORDER BY code DESC LIMIT 1"
	err := tx.QueryRowContext(ctx, query, companyID, prefix+"%").Scan(&code)

	if err != nil && err != sql.ErrNoRows {
		return code, err
	}

	if len(code) > 0 {
		runes := []rune(code)
		codeInt, err = strconv.Atoi(string(runes[8:]))
		if err != nil {
			return code, err
		}
	}

	return prefix + fmt.Sprintf("%05d", codeInt+1), nil
}
