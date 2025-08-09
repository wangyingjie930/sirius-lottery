-- -------------------------------------------------------------
-- TablePlus 6.0.0(550)
--
-- https://tableplus.com/
--
-- Database: sirius
-- Generation Time: 2025-08-09 18:18:33.3880
-- -------------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


DROP TABLE IF EXISTS `lottery_instance`;
CREATE TABLE `lottery_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `instance_id` varchar(50) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '业务活动ID',
  `instance_name` varchar(255) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '活动名称',
  `template_id` bigint unsigned NOT NULL COMMENT '关联的模板ID',
  `start_time` timestamp NOT NULL COMMENT '活动开始时间',
  `end_time` timestamp NOT NULL COMMENT '活动结束时间',
  `user_scope_json` json DEFAULT NULL COMMENT '参与用户限制',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态: 1-待上线, 2-进行中, 3-已下线',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_instance_id` (`instance_id`),
  KEY `idx_start_end_time` (`start_time`,`end_time`),
  KEY `fk_lottery_instance_template` (`template_id`),
  CONSTRAINT `fk_lottery_instance_template` FOREIGN KEY (`template_id`) REFERENCES `lottery_template` (`id`),
  CONSTRAINT `fk_lottery_pool_instance` FOREIGN KEY (`instance_id`) REFERENCES `lottery_pool` (`instance_id`),
  CONSTRAINT `fk_lottery_win_record_instance` FOREIGN KEY (`instance_id`) REFERENCES `lottery_win_record` (`instance_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs;

DROP TABLE IF EXISTS `lottery_pool`;
CREATE TABLE `lottery_pool` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `instance_id` varchar(50) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '关联的抽奖实例ID',
  `pool_name` varchar(100) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '奖池名称',
  `cost_json` json NOT NULL COMMENT '消耗的资产列表',
  `lottery_strategy` varchar(50) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '抽奖算法策略',
  `strategy_config_json` json DEFAULT NULL COMMENT '策略相关配置',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_instance_id` (`instance_id`),
  CONSTRAINT `fk_lottery_instance_pools` FOREIGN KEY (`instance_id`) REFERENCES `lottery_instance` (`instance_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs;

DROP TABLE IF EXISTS `lottery_prize`;
CREATE TABLE `lottery_prize` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `pool_id` bigint unsigned NOT NULL COMMENT '关联的奖池ID',
  `prize_id` varchar(100) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '业务奖品ID',
  `prize_name` varchar(255) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '奖品名',
  `allocated_stock` bigint NOT NULL DEFAULT '0' COMMENT '总预算库存',
  `probability` decimal(10,8) NOT NULL DEFAULT '0.00000000' COMMENT '中奖概率',
  `is_special` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否特殊奖品',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_pool_id` (`pool_id`),
  CONSTRAINT `fk_lottery_pool_prizes` FOREIGN KEY (`pool_id`) REFERENCES `lottery_pool` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs;

DROP TABLE IF EXISTS `lottery_template`;
CREATE TABLE `lottery_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `template_name` varchar(100) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '模板名称',
  `ui_style` varchar(50) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT 'UI样式标识',
  `config_json` json DEFAULT NULL COMMENT 'UI相关的配置',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态: 1-草稿, 2-已发布, 3-已归档',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs;

DROP TABLE IF EXISTS `lottery_win_record`;
CREATE TABLE `lottery_win_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` varchar(64) COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '唯一订单号',
  `instance_id` varchar(50) COLLATE utf8mb4_zh_0900_as_cs NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `prize_id` varchar(100) COLLATE utf8mb4_zh_0900_as_cs NOT NULL,
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '发放状态: 1-待发放, 2-发放成功, 3-发放失败',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_id` (`order_id`),
  KEY `idx_user_instance` (`instance_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs;

DROP TABLE IF EXISTS `transactional_messages`;
CREATE TABLE `transactional_messages` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `topic` varchar(255) COLLATE utf8mb4_zh_0900_as_cs NOT NULL,
  `key` varchar(255) COLLATE utf8mb4_zh_0900_as_cs DEFAULT NULL,
  `payload` blob NOT NULL,
  `status` varchar(20) COLLATE utf8mb4_zh_0900_as_cs NOT NULL,
  `retry_count` bigint NOT NULL DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_transactional_messages_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs;

INSERT INTO `lottery_instance` (`id`, `instance_id`, `instance_name`, `template_id`, `start_time`, `end_time`, `user_scope_json`, `status`, `created_at`, `updated_at`) VALUES
(1, 'LOTTERY_2025_SPRING', '2025年春季大抽奖', 1, '2025-03-01 00:00:00', '2025-05-31 23:59:59', '{\"user_levels\": [1, 2, 3], \"excluded_users\": [], \"min_register_days\": 7, \"max_daily_attempts\": 3}', 2, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(2, 'LOTTERY_2025_SUMMER', '2025年夏日狂欢抽奖', 2, '2025-06-01 00:00:00', '2025-08-31 23:59:59', '{\"vip_only\": false, \"user_levels\": [2, 3, 4], \"min_register_days\": 30, \"max_daily_attempts\": 5}', 1, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(3, 'LOTTERY_NEW_USER', '新用户专享抽奖', 3, '2025-01-01 00:00:00', '2025-12-31 23:59:59', '{\"new_user_only\": true, \"max_total_attempts\": 1, \"register_within_days\": 7}', 2, '2025-08-03 07:44:15', '2025-08-03 07:44:15');

INSERT INTO `lottery_pool` (`id`, `instance_id`, `pool_name`, `cost_json`, `lottery_strategy`, `strategy_config_json`, `created_at`, `updated_at`) VALUES
(1, 'LOTTERY_2025_SPRING', '春季奖池', '[{\"amount\": 100, \"asset_type\": \"points\"}, {\"amount\": 50, \"asset_type\": \"energy\"}]', 'weighted_random', '{\"algorithm\": \"mersenne_twister\", \"anti_cheating\": true, \"guarantee_mechanism\": {\"type\": \"pity\", \"threshold\": 50}}', '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(2, 'LOTTERY_2025_SUMMER', '夏日奖池', '[{\"amount\": 200, \"asset_type\": \"coins\"}]', 'progressive_odds', '{\"increment_rate\": 0.05, \"max_multiplier\": 3.0, \"base_multiplier\": 1.0, \"reset_condition\": \"daily\"}', '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(3, 'LOTTERY_NEW_USER', '新手奖池', '[{\"amount\": 1, \"asset_type\": \"free_trial\"}]', 'guaranteed_win', '{\"special_boost\": true, \"guarantee_count\": 1, \"guarantee_level\": \"medium\", \"min_prize_value\": 10}', '2025-08-03 07:44:15', '2025-08-03 08:20:27');

INSERT INTO `lottery_prize` (`id`, `pool_id`, `prize_id`, `prize_name`, `allocated_stock`, `probability`, `is_special`, `created_at`, `updated_at`) VALUES
(1, 1, 'PRIZE_SPRING_001', '一等奖-iPhone 15 Pro', 5, 0.00100000, 1, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(2, 1, 'PRIZE_SPRING_002', '二等奖-AirPods Pro', 20, 0.00500000, 1, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(3, 1, 'PRIZE_SPRING_003', '三等奖-小米手环', 100, 0.02000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(4, 1, 'PRIZE_SPRING_004', '四等奖-100元购物券', 500, 0.08000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(5, 1, 'PRIZE_SPRING_005', '五等奖-20元现金红包', 2000, 0.15000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(6, 1, 'PRIZE_SPRING_006', '安慰奖-5积分', 10000, 0.50000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(7, 1, 'PRIZE_SPRING_007', '谢谢参与', 0, 0.24400000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(8, 2, 'PRIZE_SUMMER_001', '超级大奖-特斯拉Model 3', 1, 0.00010000, 1, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(9, 2, 'PRIZE_SUMMER_002', '一等奖-MacBook Pro', 3, 0.00050000, 1, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(10, 2, 'PRIZE_SUMMER_003', '二等奖-iPad Pro', 15, 0.00300000, 1, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(11, 2, 'PRIZE_SUMMER_004', '三等奖-Switch游戏机', 50, 0.01000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(12, 2, 'PRIZE_SUMMER_005', '四等奖-500元京东卡', 200, 0.04000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(13, 2, 'PRIZE_SUMMER_006', '五等奖-100元话费', 1000, 0.10000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(14, 2, 'PRIZE_SUMMER_007', '六等奖-50元红包', 3000, 0.20000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(15, 2, 'PRIZE_SUMMER_008', '安慰奖-10积分', 8000, 0.40000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(16, 2, 'PRIZE_SUMMER_009', '谢谢参与', 0, 0.24640000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(17, 3, 'PRIZE_NEWBIE_001', '新手大礼包', 1000, 0.30000000, 1, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(18, 3, 'PRIZE_NEWBIE_002', '50元新手红包', 2000, 0.25000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(19, 3, 'PRIZE_NEWBIE_003', '会员7天体验', 3000, 0.20000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(20, 3, 'PRIZE_NEWBIE_004', '100积分奖励', 5000, 0.15000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(21, 3, 'PRIZE_NEWBIE_005', '感谢参与奖', 10000, 0.10000000, 0, '2025-08-03 07:44:15', '2025-08-03 07:44:15');

INSERT INTO `lottery_template` (`id`, `template_name`, `ui_style`, `config_json`, `status`, `created_at`, `updated_at`) VALUES
(1, '经典大转盘模板', 'lucky_wheel', '{\"colors\": [\"#FF6B6B\", \"#4ECDC4\", \"#45B7D1\", \"#96CEB4\", \"#FFEAA7\", \"#DDA0DD\"], \"animation\": {\"easing\": \"ease-out\", \"duration\": 3000}, \"background\": \"linear-gradient(135deg, #667eea 0%, #764ba2 100%)\", \"text_color\": \"#FFFFFF\", \"wheel_size\": 300, \"pointer_color\": \"#FF4757\"}', 2, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(2, '九宫格抽奖模板', 'nine_grid', '{\"animation\": {\"effect\": \"flash\", \"duration\": 200}, \"grid_size\": 3, \"background\": \"#FFFFFF\", \"text_color\": \"#212529\", \"active_color\": \"#007BFF\", \"border_color\": \"#E9ECEF\", \"item_bg_color\": \"#F8F9FA\"}', 2, '2025-08-03 07:44:15', '2025-08-03 07:44:15'),
(3, '节日主题转盘', 'lucky_wheel', '{\"colors\": [\"#E74C3C\", \"#F39C12\", \"#F1C40F\", \"#2ECC71\", \"#3498DB\", \"#9B59B6\"], \"background\": \"radial-gradient(circle, #FF6B6B 0%, #4ECDC4 100%)\", \"text_color\": \"#FFFFFF\", \"wheel_size\": 350, \"decorations\": true, \"pointer_color\": \"#C0392B\", \"festival_theme\": \"spring\"}', 2, '2025-08-03 07:44:15', '2025-08-03 07:44:15');

INSERT INTO `lottery_win_record` (`id`, `order_id`, `instance_id`, `user_id`, `prize_id`, `status`, `created_at`, `updated_at`) VALUES
(1, 'ORDER_20250801_001', 'LOTTERY_2025_SPRING', 10001, 'PRIZE_SPRING_005', 2, '2025-08-01 10:15:30', '2025-08-01 10:16:00'),
(2, 'ORDER_20250801_002', 'LOTTERY_2025_SPRING', 10002, 'PRIZE_SPRING_006', 2, '2025-08-01 11:22:15', '2025-08-01 11:22:45'),
(3, 'ORDER_20250801_003', 'LOTTERY_NEW_USER', 10003, 'PRIZE_NEWBIE_001', 2, '2025-08-01 14:30:20', '2025-08-01 14:31:10'),
(4, 'ORDER_20250801_004', 'LOTTERY_2025_SPRING', 10004, 'PRIZE_SPRING_007', 1, '2025-08-01 16:45:12', '2025-08-01 16:45:12'),
(5, 'ORDER_20250802_001', 'LOTTERY_2025_SPRING', 10005, 'PRIZE_SPRING_003', 1, '2025-08-02 09:18:33', '2025-08-02 09:18:33'),
(6, 'ORDER_20250802_002', 'LOTTERY_NEW_USER', 10006, 'PRIZE_NEWBIE_002', 2, '2025-08-02 12:55:44', '2025-08-02 12:56:20'),
(7, 'ORDER_20250802_003', 'LOTTERY_2025_SPRING', 10007, 'PRIZE_SPRING_002', 1, '2025-08-02 15:12:08', '2025-08-02 15:12:08'),
(8, 'ORDER_20250803_001', 'LOTTERY_NEW_USER', 10008, 'PRIZE_NEWBIE_004', 2, '2025-08-03 08:30:15', '2025-08-03 08:31:00'),
(9, '09b6d816-c84d-4f6f-a004-5d5ce99c21b3', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_005', 1, '2025-08-03 08:31:04', '2025-08-03 08:31:04'),
(10, '83b5145d-0348-46fd-b25b-f663367406ab', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 08:33:40', '2025-08-03 08:33:40'),
(11, '5c1a5e66-ac26-4567-941e-60b136d1f1f8', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_002', 1, '2025-08-03 08:51:51', '2025-08-03 08:51:51'),
(12, '44398b8a-b176-4c65-aaf9-c9dbc30e0af3', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_002', 1, '2025-08-03 08:52:56', '2025-08-03 08:52:56'),
(13, '2710f676-7cb3-45c5-a1f0-e659f53c4283', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_005', 1, '2025-08-03 08:53:55', '2025-08-03 08:53:55'),
(14, 'd3dfd24b-4d71-4f99-956e-ee8c405dbbbc', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 09:56:59', '2025-08-03 09:56:59'),
(15, '57c605c0-2dce-499e-94c5-e9dc4c1d86eb', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 09:58:51', '2025-08-03 09:58:51'),
(16, 'a8c12865-0d2e-4e2d-aaf4-7ae5e005b97b', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 10:08:56', '2025-08-03 10:08:56'),
(17, '14738dac-28ae-44cf-afe3-b4203ad866b3', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 10:12:17', '2025-08-03 10:12:17'),
(18, 'd8c1e2d5-47df-4311-af96-f003419eeffb', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_002', 1, '2025-08-03 10:13:35', '2025-08-03 10:13:35'),
(19, '596ba6c6-2f65-409a-a8d1-596f03d8cf8a', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 10:14:18', '2025-08-03 10:14:18'),
(20, '45822f22-b64a-4aa3-94eb-b1744758088a', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_003', 1, '2025-08-03 10:19:46', '2025-08-03 10:19:46'),
(21, '1d644f11-4694-4654-9b8f-54c5a0def707', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 10:23:20', '2025-08-03 10:23:20'),
(23, 'b8966a02-4900-42fb-b85e-28b845deace1', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_005', 1, '2025-08-03 10:39:19', '2025-08-03 10:39:19'),
(24, '4e1a2a0c-0478-4de8-be54-ebfec94b0823', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_005', 1, '2025-08-03 10:57:09', '2025-08-03 10:57:09'),
(25, '847b3ba8-4c74-4fc5-8f0b-5ea37057fdbc', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_003', 1, '2025-08-03 10:59:36', '2025-08-03 10:59:36'),
(26, 'de190726-7fc9-4edd-b7a1-ac4e63032253', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_002', 1, '2025-08-03 11:04:28', '2025-08-03 11:04:28'),
(27, '61e44eb8-6a03-4863-b77a-f58ea5736daa', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 11:08:46', '2025-08-03 11:08:46'),
(28, '03234521-6477-46b9-9e3a-1399f4f4942e', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_003', 1, '2025-08-03 11:10:59', '2025-08-03 11:10:59'),
(29, '21444d80-df45-47a0-ac4b-3fc893d44abf', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 11:20:16', '2025-08-03 11:20:16'),
(30, '30fe15a9-a0ab-4dd7-b4d2-99ecd85763d6', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_005', 1, '2025-08-03 11:27:17', '2025-08-03 11:27:17'),
(31, '0a49bc78-cb74-4b18-b7d6-59ba3c719109', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_002', 1, '2025-08-03 11:32:38', '2025-08-03 11:32:38'),
(32, '2a584cde-05b6-44b2-88c7-a8e9849d8fb4', 'LOTTERY_NEW_USER', 100, 'PRIZE_NEWBIE_004', 1, '2025-08-03 11:33:52', '2025-08-03 11:33:52');

INSERT INTO `transactional_messages` (`id`, `topic`, `key`, `payload`, `status`, `retry_count`, `created_at`, `updated_at`) VALUES
(1, 'lottery_win_events', 'e371745a-8403-4076-9611-8f09507a7399', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"09b6d816-c84d-4f6f-a004-5d5ce99c21b3\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_005\",\"status\":1}', 'PENDING', 0, '2025-08-03 16:31:04.969', '2025-08-03 16:31:04.969'),
(2, 'lottery_win_events', '9497cf26-ae9e-4b05-9364-ca161261b4b8', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"83b5145d-0348-46fd-b25b-f663367406ab\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 16:33:40.599', '2025-08-03 16:33:40.599'),
(3, 'lottery_win_events', 'e19e3e33-224a-4ae3-9dd9-bbdef676d8ab', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"5c1a5e66-ac26-4567-941e-60b136d1f1f8\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_002\",\"status\":1}', 'PENDING', 0, '2025-08-03 16:51:51.934', '2025-08-03 16:51:51.934'),
(4, 'lottery_win_events', 'd21db786-cb93-40bc-aa28-d18378574e83', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"44398b8a-b176-4c65-aaf9-c9dbc30e0af3\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_002\",\"status\":1}', 'PENDING', 0, '2025-08-03 16:52:56.245', '2025-08-03 16:52:56.245'),
(5, 'lottery_win_events', '1c86b967-b37a-4408-8e83-c8f1eea6041f', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"2710f676-7cb3-45c5-a1f0-e659f53c4283\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_005\",\"status\":1}', 'PENDING', 0, '2025-08-03 16:54:05.291', '2025-08-03 16:54:05.291'),
(6, 'lottery_win_events', '94671e0a-85bd-4b54-b25f-8d81df5e4689', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"d3dfd24b-4d71-4f99-956e-ee8c405dbbbc\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 17:56:59.121', '2025-08-03 17:56:59.121'),
(7, 'lottery_win_events', '534805a0-f578-42c7-875c-5d628e4729ea', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"57c605c0-2dce-499e-94c5-e9dc4c1d86eb\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 17:58:51.842', '2025-08-03 17:58:51.842'),
(8, 'lottery_win_events', 'f9b69cf9-7f4e-4c60-a605-2d2ab5c0769c', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"a8c12865-0d2e-4e2d-aaf4-7ae5e005b97b\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:08:56.214', '2025-08-03 18:08:56.214'),
(9, 'lottery_win_events', '831ab9a2-92d0-4e47-9861-ee962a450628', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"14738dac-28ae-44cf-afe3-b4203ad866b3\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:12:17.718', '2025-08-03 18:12:17.718'),
(10, 'lottery_win_events', '83424d9c-0c37-4efb-bb86-c9edc038011c', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"d8c1e2d5-47df-4311-af96-f003419eeffb\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_002\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:13:35.296', '2025-08-03 18:13:35.296'),
(11, 'lottery_win_events', 'f80eb81d-89cf-4295-8965-91ad83610f86', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"596ba6c6-2f65-409a-a8d1-596f03d8cf8a\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:14:18.544', '2025-08-03 18:14:18.544'),
(12, 'lottery_win_events', 'ce9b464b-df2f-4f80-92af-b7a0adff39a1', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"45822f22-b64a-4aa3-94eb-b1744758088a\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_003\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:19:46.049', '2025-08-03 18:19:46.049'),
(13, 'lottery_win_events', '691668e7-5e5b-4c72-824c-f9c1b631da91', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"1d644f11-4694-4654-9b8f-54c5a0def707\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:23:20.257', '2025-08-03 18:23:20.257'),
(14, 'lottery_win_events', '09ce9ed1-3305-46de-bfd6-bee7469260c7', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"b8966a02-4900-42fb-b85e-28b845deace1\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_005\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:39:19.749', '2025-08-03 18:39:19.749'),
(15, 'lottery_win_events', '8d6f0990-8310-46f1-9a28-2c399d88aa6f', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"4e1a2a0c-0478-4de8-be54-ebfec94b0823\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_005\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:57:09.513', '2025-08-03 18:57:09.513'),
(16, 'lottery_win_events', '1b6cb397-e485-47ef-918c-2a92f2e52287', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"847b3ba8-4c74-4fc5-8f0b-5ea37057fdbc\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_003\",\"status\":1}', 'PENDING', 0, '2025-08-03 18:59:36.202', '2025-08-03 18:59:36.202'),
(17, 'lottery_win_events', '72bdbaa8-6fe7-473c-87fd-47c7cd91ffb7', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"de190726-7fc9-4edd-b7a1-ac4e63032253\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_002\",\"status\":1}', 'PENDING', 0, '2025-08-03 19:04:28.428', '2025-08-03 19:04:28.428'),
(18, 'lottery_win_events', 'f4e425df-5bef-4a2c-8742-f96283dd6d45', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"61e44eb8-6a03-4863-b77a-f58ea5736daa\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 19:08:46.937', '2025-08-03 19:08:46.937'),
(19, 'lottery_win_events', 'b7fba3c9-af20-456b-b5db-8ac5f3706019', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"03234521-6477-46b9-9e3a-1399f4f4942e\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_003\",\"status\":1}', 'PENDING', 0, '2025-08-03 19:11:00.025', '2025-08-03 19:11:00.025'),
(20, 'lottery_win_events', '4560b73a-2b6d-4fa8-8fb5-de764e21f9a7', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"21444d80-df45-47a0-ac4b-3fc893d44abf\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 19:20:16.846', '2025-08-03 19:20:16.846'),
(21, 'lottery_win_events', 'ae794315-9ab3-471e-8a17-45015740cc46', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"30fe15a9-a0ab-4dd7-b4d2-99ecd85763d6\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_005\",\"status\":1}', 'PENDING', 0, '2025-08-03 19:27:17.097', '2025-08-03 19:27:17.097'),
(22, 'lottery_win_events', '60f092e6-55ac-4496-9f89-be9c6c3eabd7', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"0a49bc78-cb74-4b18-b7d6-59ba3c719109\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_002\",\"status\":1}', 'PENDING', 0, '2025-08-03 19:32:38.483', '2025-08-03 19:32:38.483'),
(23, 'lottery_win_events', 'b3973e90-a356-463b-bbf4-90b3d66e2b8e', '{\"id\":0,\"RequestID\":\"1111\",\"order_id\":\"2a584cde-05b6-44b2-88c7-a8e9849d8fb4\",\"instance_id\":\"LOTTERY_NEW_USER\",\"user_id\":100,\"prize_id\":\"PRIZE_NEWBIE_004\",\"status\":1}', 'PENDING', 0, '2025-08-03 19:33:52.447', '2025-08-03 19:33:52.447');



/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;