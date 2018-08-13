CREATE TABLE users (
  id           BIGSERIAL PRIMARY KEY, -- uniquely identifies the user forever
  email        TEXT NOT NULL UNIQUE, -- uniquely identifies the user at this moment, can be changed but has to remain unique
  display_name TEXT NOT NULL -- courtesy name, to be displayed, but not used to uniquely identify the user
);
