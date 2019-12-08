-- rambler up
ALTER TABLE friends ADD COLUMN "Pending" BOOLEAN;

-- rambler down
ALTER TABLE friends
DROP COLUMN "Pending";