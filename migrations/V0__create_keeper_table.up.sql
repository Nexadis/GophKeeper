CREATE TABLE IF NOT EXISTS users (
  id bigserial primary key,
  username VARCHAR(256) NOT NULL,
  hash VARCHAR(512) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  UNIQUE(username)
);

CREATE TABLE IF NOT EXISTS datas (
  id bigserial primary key,
  user_id INT NOT NULL,
  dtype VARCHAR(32) NOT NULL,
  description VARCHAR(2048),
  value VARCHAR(128) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  edited_at TIMESTAMP NOT NULL
);
