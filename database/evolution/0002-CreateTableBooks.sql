CREATE TABLE books (
  id       BIGSERIAL PRIMARY KEY,
  name     TEXT NOT NULL,
  owner_id BIGINT NOT NULL REFERENCES USERS (id),

  UNIQUE (name, owner_id)
);
