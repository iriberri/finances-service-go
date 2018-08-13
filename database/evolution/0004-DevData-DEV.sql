INSERT INTO users (email, display_name) values
  /* 1 */ ('root@foo.bar', 'I_Am_Root'),
  /* 2 */ ('david@something.net', 'David'),
  /* 3 */ ('joe@something.net', 'Joe');

INSERT INTO books (name, owner_id) values
  /* 1 */ ('David''s Book', 2);

INSERT INTO users_rights (user_id, book_id, role) values
  (1, null, 'Admin'),
  (2, 1, 'Reader');
