-- Purpose: Add a column to the folders table to store the type of folder.

ALTER TABLE folders ADD COLUMN type VARCHAR(255) NOT NULL DEFAULT 'shared';