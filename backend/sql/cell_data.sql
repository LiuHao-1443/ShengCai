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

 Date: 05/08/2024 11:44:08
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for cell_data
-- ----------------------------
DROP TABLE IF EXISTS `cell_data`;
CREATE TABLE `cell_data`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sheet_id` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `title` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `link` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `release_date` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `content` mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `abstract` varchar(1000) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `keyword` varchar(1000) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `sort_number` int(11) NULL DEFAULT NULL,
  `create_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
  `update_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0),
  `deleted` tinyint(4) UNSIGNED NULL DEFAULT 0,
  `runtime_state` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `join_sheet_id_link`(`sheet_id`, `link`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
