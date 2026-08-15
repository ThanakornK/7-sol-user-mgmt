package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUserModelDoesNotExposePasswordInJSON(t *testing.T) {
	raw, err := json.Marshal(User{Password: "bcrypt-hash"})
	if err != nil { t.Fatal(err) }
	if strings.Contains(string(raw), "bcrypt-hash") || strings.Contains(string(raw), "password") {
		t.Fatalf("password exposed: %s", raw)
	}
}
