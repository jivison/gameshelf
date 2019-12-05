-- rambler up
ALTER TABLE games ADD COLUMN "ComplexityRating" REAL;

-- rambler down
ALTER TABLE games
DROP COLUMN "ComplexityRating";