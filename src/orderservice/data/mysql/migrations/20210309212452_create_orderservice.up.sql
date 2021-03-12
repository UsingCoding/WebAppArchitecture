CREATE TABLE `order`
(
    `order_id` BINARY(16) NOT NULL,
    `ordered_at_timestamp` INTEGER NOT NULL ,
    `cost` INTEGER NOT NULL,
    PRIMARY KEY (`order_id`),
    INDEX `order_id_index` (`order_id`)
);

CREATE TABLE `menu_item`
(
    `menu_item_id` BINARY(16) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`menu_item_id`),
    INDEX `menu_item_id_index` (`menu_item_id`)
);

CREATE TABLE `order_has_menu_item`
(
    `order_has_menu_item_id` BINARY(16) NOT NULL,
    `order_id` BINARY(16) NOT NULL,
    `menu_item_id` BINARY(16) NOT NULL,
    `quantity` INTEGER,
    PRIMARY KEY (`order_has_menu_item_id`),
    INDEX `order_has_menu_item_id_index` (`order_has_menu_item_id`),
    CONSTRAINT `order_has_menu_item_menu_item_fk`
    FOREIGN KEY (`menu_item_id`) REFERENCES `menu_item`(`menu_item_id`) ON DELETE CASCADE,
    CONSTRAINT `order_has_menu_item_order_fk`
    FOREIGN KEY (`order_id`) REFERENCES `order`(`order_id`) ON DELETE CASCADE,
    UNIQUE KEY order_id_menu_item_unique_key (`order_id`, `menu_item_id`)
);