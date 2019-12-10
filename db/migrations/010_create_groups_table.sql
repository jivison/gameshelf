-- rambler up
CREATE TABLE groups (
  "ID" BIGSERIAL PRIMARY KEY,
  "Name" VARCHAR(255)
);

-- rambler down
DROP TABLE groups;