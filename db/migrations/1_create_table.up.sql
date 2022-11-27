CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  title varchar(50) NOT NULL,
  archived boolean NOT NULL
);

CREATE TABLE quizzes (
  id SERIAL PRIMARY KEY,
  category_id integer REFERENCES categories (id),
  title varchar(50) NOT NULL,
  description varchar(1000) NOT NULL,
  archived boolean NOT NULL
);

CREATE TABLE choices (
  id SERIAL PRIMARY KEY,
  quiz_id integer REFERENCES quizzes (id),
  is_correct boolean NOT NULL,
  content varchar(200) NOT NULL,
  archived boolean NOT NULL
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name varchar(30) NOT NULL,
  password_hash character(60),
  archived boolean NOT NULL
);
