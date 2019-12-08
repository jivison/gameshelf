-- rambler up
CREATE TABLE friends (
  "ID" BIGSERIAL PRIMARY KEY,
  "FrienderUsername" VARCHAR(255) REFERENCES users("Username"),
  "FriendedUsername" VARCHAR(255) REFERENCES users("Username")
);

-- rambler down
DROP TABLE friends;