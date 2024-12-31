package rapidval

import "testing"

func TestTranslator(t *testing.T) {
	tr := NewTranslator()

	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "min length error",
			err: &ValidationError{
				Field:      "name",
				MessageKey: MsgMinLength,
				MessageParams: map[string]interface{}{
					Field: "Name",
					Min:   3,
					Value: "J",
				},
			},
			expected: "Name en az 3 karakter olmalıdır",
		},
		{
			name: "required error",
			err: &ValidationError{
				Field:      "email",
				MessageKey: MsgRequired,
				MessageParams: map[string]interface{}{
					Field: "Email",
					Value: "",
				},
			},
			expected: "Email alanı zorunludur",
		},
		{
			name: "unknown message key",
			err: &ValidationError{
				Field:      "test",
				MessageKey: "unknown.key",
			},
			expected: "unknown.key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tr.Translate(tt.err)
			if got != tt.expected {
				t.Errorf("Translate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTranslatorWithCustomMessages(t *testing.T) {
	messages := map[string]string{
		MsgRequired:  "The {{.Field}} field is required",
		MsgMinLength: "The {{.Field}} must be at least {{.Min}} characters",
	}

	tr := NewTranslatorWithMessages(messages)

	err := &ValidationError{
		Field:      "name",
		MessageKey: MsgRequired,
		MessageParams: map[string]interface{}{
			Field: "Name",
		},
	}

	expected := "The Name field is required"
	got := tr.Translate(err)
	if got != expected {
		t.Errorf("Translate() = %v, want %v", got, expected)
	}
}
