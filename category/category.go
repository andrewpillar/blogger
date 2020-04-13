package category

import (
	"net/url"
	"strconv"

	"blogger/model"

	"github.com/andrewpillar/query"

	"github.com/jmoiron/sqlx"
)

type Category struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type Store struct {
	model.Store
}

var (
	_ model.Model  = (*Category)(nil)
	_ model.Loader = (*Store)(nil)
	_ model.Binder = (*Store)(nil)

	table = "categories"
)

func NewStore(db *sqlx.DB, mm ...model.Model) Store {
	s := Store{
		Store: model.Store{DB: db},
	}
	s.Bind(mm...)
	return s
}

func (c *Category) Primary() (string, int64) {
	if c == nil {
		return "id", 0
	}
	return "id", c.ID
}

func (*Category) Bind(_ ...model.Model) {}

func (c *Category) IsZero() bool {
	return c == nil || c.ID == 0 && c.Name == ""
}

func (c *Category) Values() map[string]interface{} {
	if c == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"name": c.Name,
	}
}

func (*Store) Bind(_ ...model.Model) {}

func (s Store) All(opts ...query.Option) ([]*Category, error) {
	cc := make([]*Category, 0)
	err := s.Store.All(&cc, table, opts...)
	return cc, err
}

func (s Store) Index(vals url.Values) ([]*Category, model.Paginator, error) {
	page, err := strconv.ParseInt(vals.Get("page"), 10, 64)

	if err != nil {
		page = 1
	}

	opts := []query.Option{
		model.Search("name", vals.Get("search")),
	}

	paginator, err := s.Paginate(table, page, opts...)

	if err != nil {
		return []*Category{}, paginator, nil
	}

	cc, err := s.All(append(
		opts,
		query.Limit(model.PageLimit),
		query.Offset(paginator.Offset),
	)...)
	return cc, paginator, err
}

func (s Store) Get(opts ...query.Option) (*Category, error) {
	c := &Category{}
	err := s.Store.Get(c, table, opts...)
	return c, err
}

func (s Store) Load(key string, vals []interface{}, load model.LoaderFunc) error {
	cc, err := s.All(query.Where(key, "IN", vals...))

	if err != nil {
		return err
	}

	for i := range vals {
		for _, c := range cc {
			load(i, c)
		}
	}
	return nil
}
