CREATE TABLE `users` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `email` varchar(255) UNIQUE NOT NULL,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  `role` varchar(255) NOT NULL,
  `last_modified` datetime,
  `date_created` datetime NOT NULL
);

CREATE TABLE `payment_options` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `card_type` varchar(255) NOT NULL,
  `card_number` varchar(255) NOT NULL,
  `expiry_month` int NOT NULL,
  `expiry_year` int NOT NULL,
  `name_on_card` varchar(255) NOT NULL,
  `cvv` int NOT NULL
);

CREATE TABLE `shipping_addresses` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `email_invoice` varchar(255) NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `address_line1` varchar(255) NOT NULL,
  `address_line2` varchar(255),
  `city` varchar(255) NOT NULL,
  `state` varchar(255),
  `post_code` varchar(255) NOT NULL,
  `country` varchar(255) NOT NULL
);

ALTER TABLE `payment_options` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `shipping_addresses` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

CREATE INDEX `users_index_0` ON `users` (`id`);

CREATE INDEX `users_index_1` ON `users` (`email`);

CREATE INDEX `payment_options_index_2` ON `payment_options` (`user_id`);

CREATE INDEX `shipping_addresses_index_3` ON `shipping_addresses` (`user_id`);