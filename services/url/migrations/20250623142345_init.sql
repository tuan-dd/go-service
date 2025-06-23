-- Create "shortened_urls" table

CREATE TABLE `shortened_urls` (
  `id` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL,
  `expired_at` timestamp NULL,
  `original_url` text NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
