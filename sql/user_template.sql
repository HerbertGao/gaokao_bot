SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_template
-- ----------------------------
DROP TABLE IF EXISTS `user_template`;
CREATE TABLE `user_template`  (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `user_id` int(4) NOT NULL COMMENT '用户ID',
  `template_name` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '模板名称',
  `template_content` varchar(160) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '模板内容',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_template
-- ----------------------------
INSERT INTO `user_template` VALUES (1, 0, '', '现在距离{exam}还有{time}');

SET FOREIGN_KEY_CHECKS = 1;
