-- +migrate Up
-- please read this article to understand why we use VARCHAR(191) 
-- https://www.grouparoo.com/blog/varchar-191
CREATE TABLE `permissions` (
                        `id` INT PRIMARY KEY AUTO_INCREMENT,
                        `title` VARCHAR(191) NOT NULL UNIQUE,
                        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `permissions`;