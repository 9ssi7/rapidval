package examples

import (
	"testing"
	"time"

	"github.com/9ssi7/rapidval"
)

func TestBusinessValidation(t *testing.T) {
	now := time.Now()
	business := &Business{
		Title:       "A",
		Description: "Too short",
		FoundAt:     now.Add(-24 * time.Hour),
	}

	v := rapidval.New()
	err := v.Validate(business)
	if err == nil {
		t.Error("validation should fail")
	}

	verr, ok := err.(rapidval.ValidationErrors)
	if !ok {
		t.Error("error should be ValidationErrors type")
	}

	if len(verr) != 3 {
		t.Errorf("expected 3 validation errors, got %d", len(verr))
	}

	// Hataları detaylı kontrol et
	for _, err := range verr {
		switch err.Field {
		case "Title":
			if err.MessageKey != rapidval.MsgMinLength {
				t.Errorf("unexpected Title error key: %v", err.MessageKey)
			}
			if err.MessageParams["Min"] != 3 {
				t.Errorf("unexpected min length param: %v", err.MessageParams["min"])
			}
		case "Description":
			if err.MessageKey != rapidval.MsgMinLength {
				t.Errorf("unexpected Description error key: %v", err.MessageKey)
			}
		case "FoundAt":
			if err.MessageKey != rapidval.MsgDateGreaterThan {
				t.Errorf("unexpected FoundAt error key: %v", err.MessageKey)
			}
		default:
			t.Errorf("unexpected error field: %v", err)
		}
	}
}

func TestUserValidation(t *testing.T) {
	user := &User{
		FirstName: "J",
		LastName:  "D",
		Email:     "invalid-email",
		Password:  "123",
		Age:       15,
	}

	v := rapidval.New()
	err := v.Validate(user)
	if err == nil {
		t.Error("validation should fail")
	}

	verr, ok := err.(rapidval.ValidationErrors)
	if !ok {
		t.Error("error should be ValidationErrors type")
	}

	if len(verr) != 5 { // FirstName min, LastName min, Email invalid, Password min, Age between
		t.Errorf("expected 5 validation errors, got %d", len(verr))
	}
}
