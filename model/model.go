package model

import (
	"database/sql"

	"github.com/andrewpillar/query"

	"github.com/jmoiron/sqlx"
)

type Binder interface {
	Bind(...Model)
}

type LoaderFunc func(int, Model)

type RelationFunc func(Loader, ...Model) error

type Loader interface {
	Load(string, []interface{}, LoaderFunc) error
}

type Loaders map[string]Loader

type Model interface {
	Binder

	Primary() (string, int64)

	IsZero() bool

	Values() map[string]interface{}
}

type Paginator struct {
	Next   int64
	Prev   int64
	Offset int64
	Page   int64
}

type Store struct {
	*sqlx.DB
}

type selectFunc func(interface{}, string, ...interface{}) error

var PageLimit int64 = 25

func getInt64(key string, m Model) int64 {
	if col, val := m.Primary(); key == col {
		return val
	}
	i, _ := m.Values()[key].(int64)
	return i
}

func NewLoaders() Loaders { return Loaders(make(map[string]Loader)) }

func Bind(a, b string, mm ...Model) LoaderFunc {
	return func(i int, r Model) {
		if i > len(mm) || len(mm) == 0 {
			return
		}

		m := mm[i]
		if getInt64(a, m) == getInt64(b, r) {
			m.Bind(r)
		}
	}
}

func MapKey(key string, mm ...Model) []interface{} {
	vals := make([]interface{}, 0, len(mm))

	for _, m := range mm {
		if col, val := m.Primary(); key == col {
			vals = append(vals, val)
			continue
		}

		if val, ok := m.Values()[key]; ok {
			vals = append(vals, val)
		}
	}
	return vals
}

func LoadRelations(rr map[string]RelationFunc, loaders Loaders, mm ...Model) error {
	for relation, fn := range rr {
		if err := fn(loaders.Get(relation), mm...); err != nil {
			return err
		}
	}
	return nil
}

func Relation(a, b string) RelationFunc {
	return func(l Loader, mm ...Model) error {
		return l.Load(b, MapKey(a, mm...), Bind(a, b, mm...))
	}
}

func Slice(l int, fn func(int) Model) []Model {
	mm := make([]Model, l, l)

	for i := range mm {
		mm[i] = fn(i)
	}
	return mm
}

func Where(m Model, args ...string) query.Option {
	return func(q query.Query) query.Query {
		if len(args) < 1 || m == nil || m.IsZero() {
			return q
		}

		var val interface{}

		col := args[0]

		if len(args) > 1 {
			val = m.Values()[args[1]]
		} else {
			_, val = m.Primary()
		}
		return query.Where(col, "=", val)(q)
	}
}

func Search(col, pattern string) query.Option {
	return func(q query.Query) query.Query {
		if pattern == "" {
			return q
		}
		return query.Where(col, "LIKE", "%"+pattern+"%")(q)
	}
}

func (m *Loaders) Put(name string, l Loader) {
	if (*m) == nil {
		(*m) = make(map[string]Loader)
	}
	(*m)[name] = l
}

func (m *Loaders) Get(name string) Loader { return (*m)[name] }

func (s Store) doSelect(fn selectFunc, i interface{}, table string, opts ...query.Option) error {
	opts = append([]query.Option{
		query.Columns("*"),
		query.From(table),
	}, opts...)

	q := query.Select(opts...)

	err := fn(i, q.Build(), q.Args()...)

	if err == sql.ErrNoRows {
		err = nil
	}
	return err
}

func (s Store) All(i interface{}, table string, opts ...query.Option) error {
	return s.doSelect(s.DB.Select, i, table, opts...)
}

func (s Store) Get(i interface{}, table string, opts ...query.Option) error {
	return s.doSelect(s.DB.Get, i, table, opts...)
}

func (s Store) Paginate(table string, page int64, opts ...query.Option) (Paginator, error) {
	p := Paginator{
		Page: page,
	}

	opts = append([]query.Option{
		query.Count("*"),
		query.From(table),
	}, opts...)

	q := query.Select(opts...)

	stmt, err := s.Prepare(q.Build())

	if err != nil {
		return p, err
	}

	defer stmt.Close()

	var count int64

	if err := stmt.QueryRow(q.Args()...).Scan(&count); err != nil {
		return p, err
	}

	pages := (count / PageLimit) + 1

	p.Offset = (page - 1) * PageLimit
	p.Next = page + 1
	p.Prev = page - 1

	if p.Prev < 1 {
		p.Prev = 1
	}
	if p.Next > pages {
		p.Next = pages
	}
	return p, nil
}
