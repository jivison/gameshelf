-- rambler up
ALTER TABLE games ADD COLUMN "ImgURL" VARCHAR(255);

-- rambler down
ALTER TABLE games
DROP COLUMN "ImgURL";