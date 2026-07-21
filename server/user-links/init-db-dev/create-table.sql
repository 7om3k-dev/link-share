CREATE TABLE link (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    url TEXT NOT NULL
);
