-- +migrate Up
-- MYSQL 8.0 set the role value to `user` for all old records
-- Don't change the order of enum values
-- find a better solution instead of keeping the orders!!!
ALTER TABLE `users` ADD COLUMN `role` ENUM('user','admin') NOT NULL;

-- +migrate Down
ALTER TABLE `users` Drop COLUMN `role`;