-- rambler up
CREATE TABLE group_members (
  "ID" BIGSERIAL PRIMARY KEY,
  "GroupID" BIGINT REFERENCES groups("ID"),
  "Username" VARCHAR(255) REFERENCES users("Username")
);

-- rambler down
DROP TABLE group_members;