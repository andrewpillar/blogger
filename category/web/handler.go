package web

import (
	"fmt"
	"net/http"
	"strconv"

	"blogger/category"
	"blogger/post"
	"blogger/web"

	"github.com/andrewpillar/query"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	DB *sqlx.DB

	Categories category.Store
}

func (h Handler) Index(w http.ResponseWriter, r *http.Request) {
	cc, paginator, err := h.Categories.Index(r.URL.Query())

	if err != nil {
		web.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct{
		Prev string
		Next string
		Data []*category.Category
	}{
		Prev: fmt.Sprintf("/categories?page=%d", paginator.Prev),
		Next: fmt.Sprintf("/categories?page=%d", paginator.Next),
		Data: cc,
	}
	web.JSON(w, data, http.StatusOK)
}

func (h Handler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["category"], 10, 64)

	if err != nil {
		web.JSONError(w, "Not Found", http.StatusNotFound)
		return
	}

	c, err := h.Categories.Get(query.Where("id", "=", id))

	if err != nil {
		web.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if c.IsZero() {
		web.JSONError(w, "Not Found", http.StatusNotFound)
		return
	}

	pp, paginator, err := post.NewStore(h.DB, c).Index(r.URL.Query())

	if err != nil {
		web.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct{
		Category *category.Category
		Prev     string
		Next     string
		Posts    []*post.Post
	}{
		Category: c,
		Prev:     fmt.Sprintf("/category/%d?page=%d", c.ID, paginator.Prev),
		Next:     fmt.Sprintf("/category/%d?page=%d", c.ID, paginator.Next),
		Posts:    pp,
	}
	web.JSON(w, data, http.StatusOK)
}

func RegisterRoutes(db *sqlx.DB, r *mux.Router) {
	h := Handler{
		DB:         db,
		Categories: category.NewStore(db),
	}

	r.HandleFunc("/categories", h.Index)
	r.HandleFunc("/categories/{category:[0-9]+}", h.Show)
}
