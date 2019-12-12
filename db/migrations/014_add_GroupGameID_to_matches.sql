-- rambler up
ALTER TABLE matches ADD COLUMN "GroupGameID" BIGINT REFERENCES group_games("ID");

-- rambler down
ALTER TABLE matches
DROP COLUMN "GroupGameID";