CREATE TABLE users_rights (
  user_id BIGINT NOT NULL REFERENCES USERS (id),
  book_id BIGINT REFERENCES BOOKS (id),
  role    TEXT NOT NULL,

  UNIQUE (user_id, book_id, role)
);
