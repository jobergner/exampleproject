CREATE TABLE Categories (
  ID SERIAL PRIMARY KEY,
  Title varchar(50) NOT NULL,
  Archived boolean NOT NULL
);

CREATE TABLE Quizzes (
  ID SERIAL PRIMARY KEY,
  CategoryID integer REFERENCES Categories (ID),
  Title varchar(50) NOT NULL,
  Description varchar(1000) NOT NULL,
  Archived boolean NOT NULL
);

CREATE TABLE Choices (
  ID SERIAL PRIMARY KEY,
  QuizID integer REFERENCES Quizzes (ID),
  IsCorrect boolean NOT NULL,
  Content varchar(200) NOT NULL,
  Archived boolean NOT NULL
);

CREATE TABLE Users (
  ID SERIAL PRIMARY KEY,
  Name varchar(30) NOT NULL,
  PasswordHash character(60),
  Archived boolean NOT NULL
);
