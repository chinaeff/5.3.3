package user

import "testing"

func TestRegister(t *testing.T) {
	user := User{}
	email := "test@example.com"
	password := "password123"

	result := user.Register(email, password)

	if result.Email != email {
		t.Errorf("Expected email %s, but got %s", email, result.Email)
	}

	if result.Password != password {
		t.Errorf("Expected password %s, but got %s", password, result.Password)
	}
}

func TestAuthenticate(t *testing.T) {
	user := User{Email: "test@example.com", Password: "password123"}

	if !user.Authenticate("test@example.com", "password123") {
		t.Error("Expected authentication to succeed, but it failed")
	}

	if user.Authenticate("test@example.com", "wrongpassword") {
		t.Error("Expected authentication to fail, but it succeeded")
	}

	if user.Authenticate("wrong@example.com", "password123") {
		t.Error("Expected authentication to fail, but it succeeded")
	}
}
