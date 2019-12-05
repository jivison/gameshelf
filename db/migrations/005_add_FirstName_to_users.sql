-- rambler up
ALTER TABLE users ADD COLUMN "FirstName" VARCHAR(255);

-- rambler down
ALTER TABLE users
DROP COLUMN "FirstName";