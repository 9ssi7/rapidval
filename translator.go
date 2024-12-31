package rapidval

import (
	"bytes"
	"text/template"
)

var defaultMessages = map[string]string{
	MsgRequired:        "{{.Field}} alanı zorunludur",
	MsgInvalidEmail:    "{{.Field}} geçerli bir email adresi olmalıdır",
	MsgMinLength:       "{{.Field}} en az {{.Min}} karakter olmalıdır",
	MsgMaxLength:       "{{.Field}} en fazla {{.Max}} karakter olmalıdır",
	MsgBetween:         "{{.Field}} {{.Min}} ile {{.Max}} arasında olmalıdır",
	MsgDateGreaterThan: "{{.Field}} {{.Min}} tarihinden sonra olmalıdır",
	MsgDateLessThan:    "{{.Field}} {{.Max}} tarihinden önce olmalıdır",
}

// Translator handles the translation of validation error messages.
// It uses Go's text/template package to support parameterized messages.
type Translator struct {
	messages map[string]string
	tmpl     *template.Template
}

// NewTranslator creates a new Translator with default messages.
func NewTranslator() *Translator {
	return NewTranslatorWithMessages(defaultMessages)
}

// NewTranslatorWithMessages creates a new Translator with custom messages.
// The messages map should use message keys as keys and message templates as values.
// Message templates can use Go template syntax with .Field, .Min, .Max, and .Value parameters.
func NewTranslatorWithMessages(messages map[string]string) *Translator {
	t := &Translator{
		messages: messages,
	}
	tmpl := template.New("messages")
	for key, msg := range messages {
		template.Must(tmpl.New(key).Parse(msg))
	}
	t.tmpl = tmpl
	return t
}

// Translate converts a ValidationError into a human-readable message using the configured templates.
// If the message key is not found in the templates, it returns the message key itself.
func (t *Translator) Translate(err *ValidationError) string {
	_, ok := t.messages[err.MessageKey]
	if !ok {
		return err.MessageKey
	}

	var buf bytes.Buffer
	tmpl := t.tmpl.Lookup(err.MessageKey)
	if tmpl == nil {
		return err.MessageKey
	}

	if err := tmpl.Execute(&buf, err.MessageParams); err != nil {
		valErr, ok := err.(*ValidationError)
		if ok {
			return valErr.MessageKey
		}
		return err.Error()
	}

	return buf.String()
}
