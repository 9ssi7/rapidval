package rapidval

import (
	"testing"
	"time"
)

type testStruct struct{}

func (t *testStruct) Validations() P {
	return P{}
}

type testStruct2 struct{}

func (t *testStruct2) Validations() P {
	return P{
		Required("name", ""),
		Email("email", "invalid"),
	}
}

type testStruct3 struct {
	Name  string
	Email string
	Age   int
}

func (t *testStruct3) Validations() P {
	return P{
		MinLength("Name", t.Name, 2),
		Email("Email", t.Email),
		Between("Age", t.Age, 18, 100),
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		value    interface{}
		wantErr  bool
		wantKey  string
		wantMsgs map[string]interface{}
	}{
		{
			name:    "empty string",
			field:   "name",
			value:   "",
			wantErr: true,
			wantKey: MsgRequired,
			wantMsgs: map[string]interface{}{
				"Field": "name",
				"Value": "",
			},
		},
		{
			name:    "zero int",
			field:   "age",
			value:   0,
			wantErr: true,
			wantKey: MsgRequired,
		},
		{
			name:    "nil value",
			field:   "data",
			value:   nil,
			wantErr: true,
			wantKey: MsgRequired,
		},
		{
			name:    "valid string",
			field:   "name",
			value:   "John",
			wantErr: false,
		},
		{
			name:    "valid int",
			field:   "age",
			value:   25,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Required(tt.field, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Required() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.MessageKey != tt.wantKey {
				t.Errorf("Required() message key = %v, want %v", err.MessageKey, tt.wantKey)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		value    string
		wantErr  bool
		wantKey  string
		wantMsgs map[string]interface{}
	}{
		{
			name:    "invalid email - no @",
			field:   "email",
			value:   "test.com",
			wantErr: true,
			wantKey: MsgInvalidEmail,
		},
		{
			name:    "invalid email - no domain",
			field:   "email",
			value:   "test@",
			wantErr: true,
			wantKey: MsgInvalidEmail,
		},
		{
			name:    "valid email",
			field:   "email",
			value:   "test@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Email(tt.field, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.MessageKey != tt.wantKey {
				t.Errorf("Email() message key = %v, want %v", err.MessageKey, tt.wantKey)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		value    string
		min      int
		wantErr  bool
		wantKey  string
		wantMsgs map[string]interface{}
	}{
		{
			name:    "too short",
			field:   "name",
			value:   "a",
			min:     3,
			wantErr: true,
			wantKey: MsgMinLength,
			wantMsgs: map[string]interface{}{
				"Field": "name",
				"Min":   3,
				"Value": "a",
			},
		},
		{
			name:    "exact length",
			field:   "name",
			value:   "abc",
			min:     3,
			wantErr: false,
		},
		{
			name:    "longer than min",
			field:   "name",
			value:   "abcd",
			min:     3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MinLength(tt.field, tt.value, tt.min)
			if (err != nil) != tt.wantErr {
				t.Errorf("MinLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if err.MessageKey != tt.wantKey {
					t.Errorf("MinLength() message key = %v, want %v", err.MessageKey, tt.wantKey)
				}
				if tt.wantMsgs != nil {
					for k, v := range tt.wantMsgs {
						if err.MessageParams[k] != v {
							t.Errorf("MinLength() param[%s] = %v, want %v", k, err.MessageParams[k], v)
						}
					}
				}
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   string
		max     int
		wantErr bool
		wantKey string
	}{
		{
			name:    "too long",
			field:   "name",
			value:   "abcd",
			max:     3,
			wantErr: true,
			wantKey: MsgMaxLength,
		},
		{
			name:    "exact length",
			field:   "name",
			value:   "abc",
			max:     3,
			wantErr: false,
		},
		{
			name:    "shorter than max",
			field:   "name",
			value:   "ab",
			max:     3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MaxLength(tt.field, tt.value, tt.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaxLength() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.MessageKey != tt.wantKey {
				t.Errorf("MaxLength() message key = %v, want %v", err.MessageKey, tt.wantKey)
			}
		})
	}
}

func TestBetween(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   int
		min     int
		max     int
		wantErr bool
		wantKey string
	}{
		{
			name:    "below min",
			field:   "age",
			value:   17,
			min:     18,
			max:     100,
			wantErr: true,
			wantKey: MsgBetween,
		},
		{
			name:    "above max",
			field:   "age",
			value:   101,
			min:     18,
			max:     100,
			wantErr: true,
			wantKey: MsgBetween,
		},
		{
			name:    "at min",
			field:   "age",
			value:   18,
			min:     18,
			max:     100,
			wantErr: false,
		},
		{
			name:    "at max",
			field:   "age",
			value:   100,
			min:     18,
			max:     100,
			wantErr: false,
		},
		{
			name:    "between min and max",
			field:   "age",
			value:   50,
			min:     18,
			max:     100,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Between(tt.field, tt.value, tt.min, tt.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("Between() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.MessageKey != tt.wantKey {
				t.Errorf("Between() message key = %v, want %v", err.MessageKey, tt.wantKey)
			}
		})
	}
}

func TestDateValidations(t *testing.T) {
	now := time.Now()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

	t.Run("DateGreaterThan", func(t *testing.T) {
		tests := []struct {
			name    string
			field   string
			value   time.Time
			min     time.Time
			wantErr bool
			wantKey string
		}{
			{
				name:    "past date",
				field:   "date",
				value:   past,
				min:     now,
				wantErr: true,
				wantKey: MsgDateGreaterThan,
			},
			{
				name:    "future date",
				field:   "date",
				value:   future,
				min:     now,
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := DateGreaterThan(tt.field, tt.value, tt.min)
				if (err != nil) != tt.wantErr {
					t.Errorf("DateGreaterThan() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err != nil && err.MessageKey != tt.wantKey {
					t.Errorf("DateGreaterThan() message key = %v, want %v", err.MessageKey, tt.wantKey)
				}
			})
		}
	})

	t.Run("DateLessThan", func(t *testing.T) {
		tests := []struct {
			name    string
			field   string
			value   time.Time
			max     time.Time
			wantErr bool
			wantKey string
		}{
			{
				name:    "future date",
				field:   "date",
				value:   future,
				max:     now,
				wantErr: true,
				wantKey: MsgDateLessThan,
			},
			{
				name:    "past date",
				field:   "date",
				value:   past,
				max:     now,
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := DateLessThan(tt.field, tt.value, tt.max)
				if (err != nil) != tt.wantErr {
					t.Errorf("DateLessThan() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err != nil && err.MessageKey != tt.wantKey {
					t.Errorf("DateLessThan() message key = %v, want %v", err.MessageKey, tt.wantKey)
				}
			})
		}
	})
}

func TestValidationErrors(t *testing.T) {
	errs := ValidationErrors{
		&ValidationError{
			Field:      "name",
			MessageKey: MsgRequired,
		},
		&ValidationError{
			Field:      "email",
			MessageKey: MsgInvalidEmail,
		},
	}

	errStr := errs.Error()
	if errStr == "" {
		t.Error("ValidationErrors.Error() returned empty string")
	}

	expected := "validation.required; validation.email"
	if errStr != expected {
		t.Errorf("ValidationErrors.Error() = %v, want %v", errStr, expected)
	}
}

func TestValidator(t *testing.T) {
	v := &Validator{}

	t.Run("empty params", func(t *testing.T) {
		if err := v.Validate(&testStruct{}); err != nil {
			t.Errorf("Validate() with empty params returned error: %v", err)
		}
	})

	t.Run("nil error", func(t *testing.T) {
		if err := v.Validate(&testStruct{}); err != nil {
			t.Errorf("Validate() with nil error returned error: %v", err)
		}
	})

	t.Run("with errors", func(t *testing.T) {
		err := v.Validate(&testStruct2{})
		if err == nil {
			t.Error("Validate() should return error")
		}

		verr, ok := err.(ValidationErrors)
		if !ok {
			t.Error("Validate() should return ValidationErrors")
		}

		if len(verr) != 2 {
			t.Errorf("Validate() returned %d errors, want 2", len(verr))
		}
	})
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{
			name:  "empty string",
			value: "",
			want:  true,
		},
		{
			name:  "non-empty string",
			value: "test",
			want:  false,
		},
		{
			name:  "zero int",
			value: 0,
			want:  true,
		},
		{
			name:  "non-zero int",
			value: 42,
			want:  false,
		},
		{
			name:  "false bool",
			value: false,
			want:  true,
		},
		{
			name:  "true bool",
			value: true,
			want:  false,
		},
		{
			name:  "zero time",
			value: time.Time{},
			want:  true,
		},
		{
			name:  "non-zero time",
			value: time.Now(),
			want:  false,
		},
		{
			name:  "nil",
			value: nil,
			want:  true,
		},
		{
			name:  "unsupported type",
			value: []string{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isZero(tt.value); got != tt.want {
				t.Errorf("isZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkValidations(b *testing.B) {
	b.Run("Required", func(b *testing.B) {
		field := "test"
		value := ""
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Required(field, value)
		}
	})

	b.Run("Email", func(b *testing.B) {
		field := "email"
		value := "test@example.com"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Email(field, value)
		}
	})

	b.Run("MinLength", func(b *testing.B) {
		field := "name"
		value := "test"
		min := 3
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MinLength(field, value, min)
		}
	})

	b.Run("MaxLength", func(b *testing.B) {
		field := "name"
		value := "test"
		max := 10
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MaxLength(field, value, max)
		}
	})

	b.Run("Between", func(b *testing.B) {
		field := "age"
		value := 25
		min := 18
		max := 100
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Between(field, value, min, max)
		}
	})

	b.Run("DateGreaterThan", func(b *testing.B) {
		field := "date"
		value := time.Now()
		min := value.Add(-time.Hour)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			DateGreaterThan(field, value, min)
		}
	})

	b.Run("DateLessThan", func(b *testing.B) {
		field := "date"
		value := time.Now()
		max := value.Add(time.Hour)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			DateLessThan(field, value, max)
		}
	})
}

func BenchmarkValidator(b *testing.B) {
	v := &Validator{}
	b.Run("Multiple Validations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			v.Validate(&testStruct3{
				Name:  "John",
				Email: "john@example.com",
				Age:   25,
			})
		}
	})
}

func BenchmarkTranslator(b *testing.B) {
	tr := NewTranslator()
	err := &ValidationError{
		Field:      "name",
		MessageKey: MsgMinLength,
		MessageParams: map[string]interface{}{
			Field: "Name",
			Min:   3,
			Value: "J",
		},
	}

	b.Run("Translate", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tr.Translate(err)
		}
	})
}

func BenchmarkIsZero(b *testing.B) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"string", "test"},
		{"empty string", ""},
		{"int", 42},
		{"zero int", 0},
		{"bool", true},
		{"time", time.Now()},
		{"nil", nil},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				isZero(tt.value)
			}
		})
	}
}
