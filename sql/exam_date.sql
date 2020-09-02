SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for exam_date
-- ----------------------------
DROP TABLE IF EXISTS `exam_date`;
CREATE TABLE `exam_date`  (
  `id` int(1) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `exam_year` int(4) NULL DEFAULT NULL COMMENT '考试年',
  `exam_desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '考试描述',
  `short_desc` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '考试描述（短）',
  `exam_begin_date` datetime(0) NULL DEFAULT NULL COMMENT '考试开始时间',
  `exam_end_date` datetime(0) NULL DEFAULT NULL COMMENT '考试结束时间',
  `exam_year_begin_date` datetime(0) NULL DEFAULT NULL COMMENT '考试年开始时间',
  `exam_year_end_date` datetime(0) NULL DEFAULT NULL COMMENT '考试年结束时间',
  `is_delete` tinyint(1) UNSIGNED NULL DEFAULT 0 COMMENT '是否删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of exam_date
-- ----------------------------
INSERT INTO `exam_date` VALUES (1, 2018, '2018年普通高等学校招生全国统一考试', '2018年高考', '2018-06-07 09:00:00', '2018-06-09 17:00:00', '2017-06-09 00:00:00', '2018-06-09 17:00:00', 0);
INSERT INTO `exam_date` VALUES (2, 2019, '2019年普通高等学校招生全国统一考试', '2019年高考', '2019-06-07 09:00:00', '2019-06-09 17:00:00', '2018-06-09 00:00:00', '2019-06-09 17:00:00', 0);
INSERT INTO `exam_date` VALUES (3, 2020, '2020年普通高等学校招生全国统一考试', '2020年高考', '2020-07-07 09:00:00', '2020-07-10 17:00:00', '2019-06-08 17:00:00', '2020-07-10 17:00:00', 0);
INSERT INTO `exam_date` VALUES (4, 2021, '2021年普通高等学校招生全国统一考试', '2021年高考', '2021-06-07 09:00:00', '2021-06-10 17:00:00', '2020-06-08 17:00:00', '2021-06-10 17:00:00', 0);

SET FOREIGN_KEY_CHECKS = 1;
