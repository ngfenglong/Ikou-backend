package util

import "database/sql"

func CoalesceNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}

	return ""
}
