package post

import (
	"net/url"
	"strconv"

	"blogger/category"
	"blogger/model"
	"blogger/user"

	"github.com/andrewpillar/query"

	"github.com/jmoiron/sqlx"
)

type Post struct {
	ID         int64  `db:"id"`
	UserID     int64  `db:"user_id"`
	CategoryID int64  `db:"category_id"`
	Title      string `db:"title"`
	Body       string `db:"body"`

	User     *user.User         `db:"-"`
	Category *category.Category `db:"-"`
	Tags     []*Tag             `db:"-"`
}

type Store struct {
	model.Store

	User     *user.User
	Category *category.Category
}

var (
	_ model.Model  = (*Post)(nil)
	_ model.Loader = (*Store)(nil)
	_ model.Binder = (*Store)(nil)

	table     = "posts"
	relations = map[string]model.RelationFunc{
		"user":     model.Relation("user_id", "id"),
		"category": model.Relation("category_id", "id"),
		"tag":      model.Relation("id", "post_id"),
	}
)

func NewStore(db *sqlx.DB, mm ...model.Model) Store {
	s := Store{
		Store: model.Store{DB: db},
	}
	s.Bind(mm...)
	return s
}

func Model(pp []*Post) func(int) model.Model {
	return func(i int) model.Model {
		return pp[i]
	}
}

func LoadRelations(loaders model.Loaders, pp ...*Post) error {
	mm := model.Slice(len(pp), Model(pp))
	return model.LoadRelations(relations, loaders, mm...)
}

func (p *Post) Primary() (string, int64) {
	if p == nil {
		return "id", 0
	}
	return "id", p.ID
}

func (p *Post) Bind(mm ...model.Model) {
	for _, m := range mm {
		switch m.(type) {
		case *user.User:
			p.User = m.(*user.User)
		case *category.Category:
			p.Category = m.(*category.Category)
		case *Tag:
			p.Tags = append(p.Tags, m.(*Tag))
		}
	}
}

func (p *Post) IsZero() bool {
	return p == nil || p.ID == 00 && p.UserID == 0 && p.CategoryID == 0 && p.Title == "" && p.Body == ""
}

func (p *Post) Values() map[string]interface{} {
	if p == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"user_id":     p.UserID,
		"category_id": p.CategoryID,
		"title":       p.Title,
		"body":        p.Body,
	}
}

func (s *Store) Bind(mm ...model.Model) {
	for _, m := range mm {
		switch m.(type) {
		case *user.User:
			s.User = m.(*user.User)
		case *category.Category:
			s.Category = m.(*category.Category)
		}
	}
}

func (s Store) All(opts ...query.Option) ([]*Post, error) {
	pp := make([]*Post, 0)

	opts = append([]query.Option{
		model.Where(s.Category, "category_id"),
	}, opts...)

	err := s.Store.All(&pp, table, opts...)
	return pp, err
}

func (s Store) Index(vals url.Values) ([]*Post, model.Paginator, error) {
	page, err := strconv.ParseInt(vals.Get("page"), 10, 64)

	if err != nil {
		page = 1
	}

	opts := []query.Option{
		model.Search("title", vals.Get("search")),
		WhereTag(vals.Get("tag")),
	}

	paginator, err := s.Paginate(table, page, opts...)

	if err != nil {
		return []*Post{}, paginator, err
	}

	pp, err := s.All(append(
		opts,
		query.Limit(model.PageLimit),
		query.Offset(paginator.Offset),
	)...)
	return pp, paginator, err
}

func (s Store) Get(opts ...query.Option) (*Post, error) {
	p := &Post{
		User:     s.User,
		Category: s.Category,
	}
	err := s.Store.Get(p, table, opts...)
	return p, err
}

func (s Store) Load(key string, vals []interface{}, load model.LoaderFunc) error {
	pp, err := s.All(query.Where(key, "IN", vals...))

	if err != nil {
		return err
	}

	for i := range vals {
		for _, p := range pp {
			load(i, p)
		}
	}
	return nil
}
