CREATE TABLE IF NOT EXISTS users (
  id   INTEGER     PRIMARY KEY,
  user_name text   NOT NULL UNIQUE,
  password  text    NOT NULL
);

CREATE TABLE IF NOT EXISTS calendar_events (
  id   INTEGER     PRIMARY KEY,
  title text       NOT NULL,
  date_time datetime    NOT NULL,
  owner_id  INTEGER  NOT NULL
);

CREATE TABLE IF NOT EXISTS participations (
  id   INTEGER     PRIMARY KEY,
  user_id INTEGER  NOT NULL,
  event_id INTEGER  NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  id   text     PRIMARY KEY,
  user_id  INTEGER  NOT NULL
);
