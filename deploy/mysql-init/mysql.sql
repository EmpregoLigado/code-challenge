--MYSQL
CREATE DATABASE IF NOT EXISTS jobdb;
USE jobdb;
CREATE TABLE IF NOT EXISTS jobs (
    partner_id int(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    category_id int(64) NOT NULL,
    `status` VARCHAR(255) NOT NULL,
    expires_at DATETIME NOT NULL,
    PRIMARY KEY(partner_id)
);
