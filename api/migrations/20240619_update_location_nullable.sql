-- Use ALTER TABLE to modify the location column to be nullable
ALTER TABLE person
ALTER COLUMN location DROP NOT NULL;
