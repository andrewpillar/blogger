package user

import (
	"blogger/model"

	"github.com/andrewpillar/query"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type Store struct {
	model.Store
}

var (
	_ model.Model  = (*User)(nil)
	_ model.Loader = (*Store)(nil)
	_ model.Binder = (*Store)(nil)

	table = "users"
)

func NewStore(db *sqlx.DB, mm ...model.Model) Store {
	s := Store{
		Store: model.Store{DB: db},
	}
	s.Bind(mm...)
	return s
}

func (u *User) Primary() (string, int64) {
	if u == nil {
		return "id", 0
	}
	return "id", u.ID
}

func (*User) Bind(_ ...model.Model) {}

func (u *User) IsZero() bool {
	return u == nil || u.ID == 0 && u.Email == "" && u.Username == "" && u.Password == ""
}

func (u *User) Values() map[string]interface{} {
	if u == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"email":    u.Email,
		"username": u.Username,
		"password": u.Password,
	}
}

func (*Store) Bind(_ ...model.Model) {}

func (s Store) All(opts ...query.Option) ([]*User, error) {
	uu := make([]*User, 0)
	err := s.Store.All(&uu, table, opts...)
	return uu, err
}

func (s Store) Get(opts ...query.Option) (*User, error) {
	u := &User{}
	err := s.Store.Get(u, table, opts...)
	return u, err
}

func (s Store) Load(key string, vals []interface{}, load model.LoaderFunc) error {
	uu, err := s.All(query.Where(key, "IN", vals...))

	if err != nil {
		return err
	}

	for i := range vals {
		for _, u := range uu {
			load(i, u)
		}
	}
	return nil
}
