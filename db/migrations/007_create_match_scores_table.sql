-- rambler up
CREATE TABLE match_scores (
  "ID" BIGSERIAL PRIMARY KEY,
  "MatchID" BIGINT REFERENCES matches("ID"),
  "GameID" BIGINT REFERENCES games("ID"),
  "PlayerUserName" VARCHAR(255) REFERENCES users("Username"),
  "PlayerDisplayName" VARCHAR(255),
  "BaseScore" REAL,
  "IsWinner" BOOLEAN,
  "FinalScore" REAL
);

-- rambler down
DROP TABLE match_scores;