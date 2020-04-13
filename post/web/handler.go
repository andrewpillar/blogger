package web

import (
	"fmt"
	"net/http"
	"strconv"

	"blogger/category"
	"blogger/model"
	"blogger/post"
	"blogger/user"
	"blogger/web"

	"github.com/andrewpillar/query"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	DB      *sqlx.DB
	Posts   post.Store
	Loaders model.Loaders
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) {
	pp, paginator, err := h.Posts.Index(r.URL.Query())

	if err != nil {
		println(err.Error())
		web.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := post.LoadRelations(h.Loaders, pp...); err != nil {
		web.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct{
		Prev string
		Next string
		Data []*post.Post
	}{
		Prev: fmt.Sprintf("/posts?page=%d", paginator.Prev),
		Next: fmt.Sprintf("/posts?page=%d", paginator.Next),
		Data: pp,
	}
	web.JSON(w, data, http.StatusOK)
}

func (h Handler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["post"], 10, 64)

	if err != nil {
		web.JSONError(w, "Not Found", http.StatusNotFound)
		return
	}

	p, err := h.Posts.Get(query.Where("id", "=", id))

	if err != nil {
		web.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.IsZero() {
		web.JSONError(w, "Not Found", http.StatusNotFound)
		return
	}

	if err := post.LoadRelations(h.Loaders, p); err != nil {
		web.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	web.JSON(w, p, http.StatusOK)
}

func RegisterRoutes(db *sqlx.DB, r *mux.Router) {
	loaders := model.NewLoaders()
	loaders.Put("tag", post.NewTagStore(db))
	loaders.Put("user", user.NewStore(db))
	loaders.Put("category", category.NewStore(db))

	h := Handler{
		DB:      db,
		Posts:   post.NewStore(db),
		Loaders: loaders,
	}

	r.HandleFunc("/posts", h.Index)
	r.HandleFunc("/posts/{post:[0-9]+}", h.Show)
}
