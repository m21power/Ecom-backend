CREATE TABLE IF NOT EXISTS user(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `firstName` VARCHAR(255) NOT NULL,
    `lastName` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(`id`),
    UNIQUE(`email`)
)

CREATE TABLE IF NOT EXISTS `products`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `imgage` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `quantity` INT NOT NULL,
    `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(`id`)
)

CREATE TABLE IF NOT EXISTS `orders`(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `userId` INT UNSIGNED NOT NULL,
    `total` DECIMAL(10, 2) NOT NULL,
    `status` ENUM('pending', 'completed', 'cancelled') NOT NULL DEFAULT 'pending',
    `address` TEXT NOT NULL,   
    `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(`id`),
    FOREIGN KEY(`userId`) REFERENCES `user`(`id`)
)