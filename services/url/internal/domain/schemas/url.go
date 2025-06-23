package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	entOrm "github.com/tuan-dd/go-pkg/orm/ent"
)

type ShortenedURL struct {
	ent.Schema
}

func (ShortenedURL) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.Time("expired_at").Optional(),
		field.String("original_url").SchemaType(map[string]string{
			dialect.MySQL: "TEXT",
		}),
	}
}

func (ShortenedURL) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "shortened_urls"},
	}
}

func (ShortenedURL) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entOrm.TimeMixin{},
	}
}
