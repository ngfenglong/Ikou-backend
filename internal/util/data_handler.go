package util

import "database/sql"

func CoalesceNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func CoalesceNullInt(i sql.NullInt32) int {
	if i.Valid {
		return int(i.Int32)
	}

	return 0
}
