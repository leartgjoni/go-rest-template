CREATE TABLE articles(
                         id serial PRIMARY KEY,
                         slug VARCHAR (255) UNIQUE NOT NULL,
                         title VARCHAR (255) NOT NULL,
                         body TEXT,
                         user_id INTEGER REFERENCES users(id) NOT NULL,
                         created_at TIMESTAMPTZ NOT NULL,
                         updated_at TIMESTAMPTZ
);