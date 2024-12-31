package examples

import (
	"bytes"
	"text/template"

	"github.com/9ssi7/rapidval"
)

var messages = map[string]string{
	rapidval.MsgRequired:        "{{.Field}} alanı zorunludur",
	rapidval.MsgInvalidEmail:    "{{.Field}} geçerli bir email adresi olmalıdır",
	rapidval.MsgMinLength:       "{{.Field}} en az {{.Min}} karakter olmalıdır",
	rapidval.MsgMaxLength:       "{{.Field}} en fazla {{.Max}} karakter olmalıdır",
	rapidval.MsgBetween:         "{{.Field}} {{.Min}} ile {{.Max}} arasında olmalıdır",
	rapidval.MsgDateGreaterThan: "{{.Field}} {{.Min}} tarihinden sonra olmalıdır",
	rapidval.MsgDateLessThan:    "{{.Field}} {{.Max}} tarihinden önce olmalıdır",
}

type Translator struct {
	messages map[string]string
	tmpl     *template.Template
}

func NewTranslator() *Translator {
	t := &Translator{
		messages: messages,
	}
	// Her mesajı template olarak parse et
	tmpl := template.New("messages")
	for key, msg := range messages {
		template.Must(tmpl.New(key).Parse(msg))
	}
	t.tmpl = tmpl
	return t
}

func (t *Translator) Translate(err *rapidval.ValidationError) string {
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
		valErr, ok := err.(*rapidval.ValidationError)
		if ok {
			return valErr.MessageKey
		}
		return err.Error()
	}

	return buf.String()
}
