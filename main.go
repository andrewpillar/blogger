package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	categoryweb "blogger/category/web"
	postweb "blogger/post/web"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := mux.NewRouter()

	db, err := sqlx.Open("sqlite3", "db.sqlite")

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	postweb.RegisterRoutes(db, router)
	categoryweb.RegisterRoutes(db, router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("listening on", srv.Addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	sig := <-c

	srv.Shutdown(ctx)
	fmt.Println("received signal", sig, "shutting down")
}
