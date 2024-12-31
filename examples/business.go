package examples

import (
	"time"

	"github.com/9ssi7/rapidval"
)

type Business struct {
	Title       string
	Description string
	FoundAt     time.Time
}

func (b *Business) Validate(v *rapidval.Validator) error {
	now := time.Now()
	return v.Validate(rapidval.P{
		rapidval.Required("Title", b.Title),
		rapidval.Required("Description", b.Description),
		rapidval.Required("FoundAt", b.FoundAt),
		rapidval.MinLength("Title", b.Title, 3),
		rapidval.MaxLength("Title", b.Title, 100),
		rapidval.MinLength("Description", b.Description, 10),
		rapidval.MaxLength("Description", b.Description, 1000),
		rapidval.DateGreaterThan("FoundAt", b.FoundAt, now),
		rapidval.DateLessThan("FoundAt", b.FoundAt, now.Add(24*time.Hour)),
	})
}
