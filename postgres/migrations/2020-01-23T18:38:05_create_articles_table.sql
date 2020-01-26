CREATE TABLE articles(
                         id serial PRIMARY KEY,
                         slug VARCHAR (255) UNIQUE NOT NULL,
                         title VARCHAR (255) NOT NULL,
                         body TEXT,
                         user_id INTEGER REFERENCES users(id),
                         created_at TIMESTAMP NOT NULL,
                         updated_at TIMESTAMP
);