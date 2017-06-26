package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Needed for postgres DB
	"github.com/mstojcevich/lambda-ng-go/config"
)

// SQL to create all of the required tables for Lambda
var schema = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username varchar NOT NULL UNIQUE,
    password varchar NOT NULL,
    creation_date timestamp NOT NULL,
    api_key varchar NOT NULL UNIQUE,
    encryption_enabled boolean NOT NULL,
    theme_name varchar
);

CREATE TABLE IF NOT EXISTS authorities (
	id SERIAL PRIMARY KEY,
	user_id integer NOT NULL UNIQUE,
	authority_level integer NOT NULL
);

CREATE TABLE IF NOT EXISTS files (
	id SERIAL PRIMARY KEY,
	owner integer NOT NULL,
	name varchar NOT NULL UNIQUE,
	extension varchar NOT NULL,
	encrypted boolean NOT NULL,
	local_name varchar,
	upload_date timestamp,
	has_thumbnail boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS pastes (
	id SERIAL PRIMARY KEY,
	owner integer,
	name varchar NOT NULL UNIQUE,
	content_json varchar NOT NULL,
	is_code boolean NOT NULL,
	upload_date timestamp
);

CREATE TABLE IF NOT EXISTS thumbnails (
	id SERIAL PRIMARY KEY,
	parent_name varchar NOT NULL,
	width integer NOT NULL,
	height integer NOT NULL,
	url varchar NOT NULL UNIQUE
);
`

// DB is a connection to the primary Lambda database
var DB = initDatabase(config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)

func initDatabase(host string, port uint16, user string, pass string, dbName string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d "+
		"user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbName))

	// Error establishing connection to DB
	if err != nil {
		panic(err)
	}

	// Create the tables if they don't exist
	db.Exec(schema)

	return db
}
