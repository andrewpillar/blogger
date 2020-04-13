#!/bin/sh

set -e

for bin in go sqlite3; do
	if ! hash "$bin"; then
		>&2 printf "missing binary: %s\n" "$bin"
		exit 1
	fi
done

db() {
	rm -f db.sqlite

	cat schema.sql | sqlite3 db.sqlite

	while read -r row; do
		echo -n "$row" | awk -F"," '{
			_email    = "\x27" $1 "\x27"
			_username = "\x27" $2 "\x27"
			_password = "\x27" $3 "\x27"

			print "INSERT INTO users (email, username, password)"
			print "VALUES (" _email ", " _username ", " _password ")"
		}' | sqlite3 db.sqlite
	done <<< $(cat _data/users.csv)

	while read -r row; do
		printf "INSERT INTO categories (name) VALUES ('%s')" "$row" | sqlite3 db.sqlite
	done <<< $(cat _data/categories.csv)

	while read -r row; do
		echo -n "$row" | awk -F"," '{
			_user_id     = "\x27" $1 "\x27"
			_category_id = "\x27" $2 "\x27"
			_title       = "\x27" $3 "\x27"
			_body        = "\x27" $4 "\x27"

			print "INSERT INTO posts (user_id, category_id, title, body)"
			print "VALUES (" _user_id ", " _category_id ", " _title ", " _body ")"
		}' | sqlite3 db.sqlite
	done <<< $(cat _data/posts.csv)

	while read -r row; do
		echo -n "$row" | awk -F"," '{
			_post_id = "\x27" $1 "\x27"
			_name    = "\x27" $2 "\x27"

			print "INSERT INTO post_tags (post_id, name)"
			print "VALUES (" _post_id ", " _name ")"
		}' | sqlite3 db.sqlite
	done <<< $(cat _data/tags.csv)
}

db
go build -o blogger.out
