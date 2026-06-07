package util_test

import (
	"database/sql"
	"testing"

	"github.com/ngfenglong/ikou-backend/internal/util"
)

func TestCoalesceNullString_Valid_ReturnsString(t *testing.T) {
	ns := sql.NullString{String: "hello", Valid: true}
	if got := util.CoalesceNullString(ns); got != "hello" {
		t.Errorf("want hello, got %q", got)
	}
}

func TestCoalesceNullString_Invalid_ReturnsEmpty(t *testing.T) {
	ns := sql.NullString{String: "ignored", Valid: false}
	if got := util.CoalesceNullString(ns); got != "" {
		t.Errorf("want empty string, got %q", got)
	}
}

func TestCoalesceNullInt_Valid_ReturnsInt(t *testing.T) {
	ni := sql.NullInt32{Int32: 42, Valid: true}
	if got := util.CoalesceNullInt(ni); got != 42 {
		t.Errorf("want 42, got %d", got)
	}
}

func TestCoalesceNullInt_Invalid_ReturnsZero(t *testing.T) {
	ni := sql.NullInt32{Int32: 99, Valid: false}
	if got := util.CoalesceNullInt(ni); got != 0 {
		t.Errorf("want 0, got %d", got)
	}
}
