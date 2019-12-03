-- rambler up
CREATE TABLE matches (
    "ID" BIGSERIAL PRIMARY KEY,
    "GameID" BIGINT REFERENCES games("ID"),
    "DatePlayed" DATE,
    "HostUserName" VARCHAR(255) REFERENCES users("Username")
);

-- rambler down
DROP TABLE matches;