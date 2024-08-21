CREATE TABLE IF NOT EXISTS users (

                                     id INT NOT NULL UNIQUE PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     email TEXT NOT NULL UNIQUE
);