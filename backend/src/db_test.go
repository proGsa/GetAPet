package main

import (
	"os"
	"testing"
)

func TestDBConnectionSuccess(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("DB should be reachable, got error: %v", err)
	}
}

func TestDBConnectionError(t *testing.T) {
	// Подменяем переменные окружения на неверные
	os.Setenv("POSTGRES_HOST", "invalid_host")

	db, err := getDB()
	if err != nil {
		t.Fatalf("expected connection object, got error: %v", err)
	}

	// Ping должен упасть
	err = db.Ping()
	if err == nil {
		t.Fatalf("expected error when connecting to invalid host, got nil")
	}
}