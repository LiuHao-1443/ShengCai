/*
 Navicat Premium Data Transfer

 Source Server         : localhost-mysql
 Source Server Type    : MySQL
 Source Server Version : 50744
 Source Host           : localhost:3336
 Source Schema         : shengcai

 Target Server Type    : MySQL
 Target Server Version : 50744
 File Encoding         : 65001

 Date: 05/08/2024 11:44:14
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sheet_info
-- ----------------------------
DROP TABLE IF EXISTS `sheet_info`;
CREATE TABLE `sheet_info`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sheet_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `sheet_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `update_log` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `create_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
  `update_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0),
  `deleted` tinyint(4) UNSIGNED NULL DEFAULT 0,
  `runtime_state` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `sheet_id`(`sheet_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
