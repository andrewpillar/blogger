

CREATE TABLE users (
	id       INTEGER PRIMARY KEY AUTOINCREMENT,
	email    TEXT,
	username TEXT,
	password TEXT
);

CREATE TABLE categories (
	id   INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT
);

CREATE TABLE posts (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id     INT,
	category_id INT,
	title       TEXT,
	body        TEXT
);

CREATE TABLE post_tags (
	id      INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INT,
	name    TEXT
);
