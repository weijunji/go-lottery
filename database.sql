CREATE DATABASE IF NOT EXISTS lottery_test;

USE `lottery_test`;

CREATE TABLE `lotteries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(32) NOT NULL,
  `description` text,
  `permanent` bigint unsigned NOT NULL,
  `temporary` bigint unsigned NOT NULL,
  `start_time` datetime(3) NOT NULL,
  `end_time` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `award_infos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `lottery` bigint unsigned NOT NULL,
  `name` varchar(32) DEFAULT NULL,
  `type` bigint unsigned DEFAULT NULL,
  `description` text DEFAULT NULL,
  `pic` text DEFAULT NULL,
  `total` bigint DEFAULT NULL,
  `display_rate` bigint DEFAULT NULL,
  `rate` bigint NOT NULL,
  `value` int unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `fk_award_infos_fkey` (`lottery`),
  CONSTRAINT `fk_award_infos_fkey` FOREIGN KEY (`lottery`) REFERENCES `lotteries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `awards` (
  `award` bigint unsigned NOT NULL AUTO_INCREMENT,
  `lottery` bigint unsigned NOT NULL,
  `remain` bigint NOT NULL,
  PRIMARY KEY (`award`),
  KEY `idx_awards_lottery` (`lottery`),
  CONSTRAINT `fk_awards_fkey1` FOREIGN KEY (`lottery`) REFERENCES `lotteries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_awards_fkey2` FOREIGN KEY (`award`) REFERENCES `award_infos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `access_token` varchar(128) NOT NULL,
  `token_type` bigint NOT NULL,
  `role` bigint DEFAULT 1,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `access_token` (`access_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `winning_infos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user` bigint unsigned NOT NULL,
  `award` bigint unsigned NOT NULL,
  `lottery` bigint unsigned NOT NULL,
  `address` tinytext,
  `handout` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_winning_infos_user` (`user`),
  KEY `idx_winning_infos_award` (`award`),
  KEY `idx_winning_infos_lottery` (`lottery`),
  CONSTRAINT `fk_winning_infos_fkey1` FOREIGN KEY (`user`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_winning_infos_fkey2` FOREIGN KEY (`award`) REFERENCES `award_infos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_winning_infos_fkey3` FOREIGN KEY (`lottery`) REFERENCES `lotteries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
