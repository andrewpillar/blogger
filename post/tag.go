package post

import (
	"blogger/model"

	"github.com/andrewpillar/query"

	"github.com/jmoiron/sqlx"
)

type Tag struct {
	ID     int64  `db:"id"`
	PostID int64  `db:"post_id"`
	Name   string `db:"name"`

	Post *Post `db:"-"`
}

type TagStore struct {
	model.Store

	Post *Post
}

var (
	_ model.Model  = (*Tag)(nil)
	_ model.Loader = (*TagStore)(nil)
	_ model.Binder = (*TagStore)(nil)

	tagTable = "post_tags"
)

func NewTagStore(db *sqlx.DB, mm ...model.Model) TagStore {
	s := TagStore{
		Store: model.Store{DB: db},
	}
	s.Bind(mm...)
	return s
}

func WhereTag(tag string) query.Option {
	return func(q query.Query) query.Query {
		if tag == "" {
			return q
		}
		return query.WhereQuery("id", "IN",
			query.Select(
				query.Columns("post_id"),
				query.From(tagTable),
				query.Where("name", "=", tag),
			),
		)(q)
	}
}

func (t *Tag) Primary() (string, int64) {
	if t == nil {
		return "id", 0
	}
	return "id", t.ID
}

func (t *Tag) Bind(mm ...model.Model) {
	for _, m := range mm {
		switch m.(type) {
		case *Post:
			t.Post = m.(*Post)
		}
	}
}

func (t *Tag) IsZero() bool {
	return t == nil || t.ID == 0 && t.PostID == 0 && t.Name == ""
}

func (t *Tag) Values() map[string]interface{} {
	if t == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"post_id": t.PostID,
		"name":    t.Name,
	}
}

func (s *TagStore) Bind(mm ...model.Model) {
	for _, m := range mm {
		switch m.(type) {
		case *Post:
			s.Post = m.(*Post)
		}
	}
}

func (s TagStore) All(opts ...query.Option) ([]*Tag, error) {
	tt := make([]*Tag, 0)
	err := s.Store.All(&tt, tagTable, opts...)
	return tt, err
}

func (s TagStore) Get(opts ...query.Option) (*Tag, error) {
	t := &Tag{
		Post: s.Post,
	}
	err := s.Store.Get(t, tagTable, opts...)
	return t, err
}

func (s TagStore) Load(key string, vals []interface{}, load model.LoaderFunc) error {
	tt, err := s.All(query.Where(key, "IN", vals...))

	if err != nil {
		return err
	}

	for i := range vals {
		for _, t := range tt {
			load(i, t)
		}
	}
	return nil
}
