CREATE TABLE users(
                      id serial PRIMARY KEY,
                      username VARCHAR (50) NOT NULL,
                      email VARCHAR (255) UNIQUE NOT NULL,
                      password VARCHAR (255) NOT NULL,
                      created_at TIMESTAMPTZ NOT NULL,
                      updated_at TIMESTAMPTZ
);