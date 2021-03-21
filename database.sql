CREATE DATABASE lottery_test;

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
