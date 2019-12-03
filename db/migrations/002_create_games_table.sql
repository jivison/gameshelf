-- rambler up
CREATE TABLE games (
  "ID" BIGSERIAL PRIMARY KEY,
  "Title" VARCHAR(255) UNIQUE NOT NULL,
  "Year" INTEGER,
  "BggID" INTEGER,
  "user_name" VARCHAR(255) REFERENCES users("Username")
);

-- rambler down
DROP TABLE games;