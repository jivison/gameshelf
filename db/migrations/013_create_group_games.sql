-- rambler up
CREATE TABLE group_games (
  "ID" BIGSERIAL PRIMARY KEY,
  "GroupID" BIGINT REFERENCES groups("ID"),
  "GameID" BIGINT REFERENCES games("ID"),
  "TimesPlayed" INT
);

-- rambler down
DROP TABLE group_games;