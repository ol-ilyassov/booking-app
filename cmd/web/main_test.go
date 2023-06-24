package main

import "testing"

func TestRun(t *testing.T) {
	// db, err := run()
	_, err := run()
	if err != nil {
		t.Error("failed run()")
	}
	// defer db.SQL.Close()
}
