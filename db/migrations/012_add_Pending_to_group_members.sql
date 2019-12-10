-- rambler up
ALTER TABLE group_members ADD COLUMN "Pending" BOOLEAN;

-- rambler down
ALTER TABLE group_members
DROP COLUMN "Pending";