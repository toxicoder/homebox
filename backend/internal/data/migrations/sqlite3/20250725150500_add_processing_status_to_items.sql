-- +goose Up
ALTER TABLE `items` ADD COLUMN `processing_status` TEXT NOT NULL DEFAULT 'local_only';

-- +goose Down
ALTER TABLE `items` DROP COLUMN `processing_status`;
