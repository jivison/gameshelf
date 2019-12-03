-- rambler up
CREATE TABLE users (
  "Username" VARCHAR(255) PRIMARY KEY,
  "HashedPassword" VARCHAR(255)
);

-- rambler down
DROP TABLE users;