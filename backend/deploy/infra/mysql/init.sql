CREATE DATABASE IF NOT EXISTS `kubehostwarden` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `kubehostwarden`.`host`;

CREATE TABLE `kubehostwarden`.`host` (
    `id` VARCHAR(255),
    `hostname` VARCHAR(64),
    `os` VARCHAR(64),
    `os_version` VARCHAR(64),
    `kernel` VARCHAR(64),
    `kernel_version` VARCHAR(64),
    `arch` VARCHAR(64),
    `ip_addr` VARCHAR(64),
    `memory_total` VARCHAR(64),
    `disk_total` VARCHAR(64),
    `owner_id` VARCHAR(255),
    `owner` VARCHAR(64),
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `kubehostwarden`.`user`;

CREATE TABLE `kubehostwarden`.`user` (
    `id` VARCHAR(255),
    `username` VARCHAR(64),
    `password` VARCHAR(64),
    `email` VARCHAR(64),
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `kubehostwarden`.`threshold_info`;

CREATE TABLE `kubehostwarden`.`threshold_info` (
    `id` VARCHAR(255),
    `host_id` VARCHAR(255),
    `metric` VARCHAR(64),
    `sub_metric` VARCHAR(64),
    `threshold` FLOAT,
    `type` VARCHAR(64),
    `entry_id` INT,
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
     PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
