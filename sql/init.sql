SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for exam_date
-- ----------------------------
DROP TABLE IF EXISTS `exam_date`;
CREATE TABLE `exam_date` (
  `id` int(1) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `exam_year` int(4) DEFAULT NULL COMMENT '考试年',
  `exam_desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '考试描述',
  `short_desc` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '考试描述（短）',
  `exam_begin_date` datetime DEFAULT NULL COMMENT '考试开始时间',
  `exam_end_date` datetime DEFAULT NULL COMMENT '考试结束时间',
  `exam_year_begin_date` datetime DEFAULT NULL COMMENT '考试年开始时间',
  `exam_year_end_date` datetime DEFAULT NULL COMMENT '考试年结束时间',
  `is_delete` tinyint(1) unsigned DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of exam_date
-- ----------------------------
BEGIN;
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (1, 2018, '2018年普通高等学校招生全国统一考试', '2018年高考', '2018-06-07 09:00:00', '2018-06-09 17:00:00', '2017-06-09 00:00:00', '2018-06-09 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (2, 2019, '2019年普通高等学校招生全国统一考试', '2019年高考', '2019-06-07 09:00:00', '2019-06-09 17:00:00', '2018-06-09 00:00:00', '2019-06-09 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (3, 2020, '2020年普通高等学校招生全国统一考试', '2020年高考', '2020-07-07 09:00:00', '2020-07-10 17:00:00', '2019-06-08 17:00:00', '2020-07-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (4, 2021, '2021年普通高等学校招生全国统一考试', '2021年高考', '2021-06-07 09:00:00', '2021-06-10 17:00:00', '2020-06-08 17:00:00', '2021-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (5, 2022, '2022年普通高等学校招生全国统一考试', '2022年高考', '2022-06-07 09:00:00', '2022-06-10 17:00:00', '2021-06-08 17:00:00', '2022-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (6, 2023, '2023年普通高等学校招生全国统一考试', '2023年高考', '2023-06-07 09:00:00', '2023-06-10 17:00:00', '2022-07-09 17:00:00', '2023-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (7, 2024, '2024年普通高等学校招生全国统一考试', '2024年高考', '2024-06-07 09:00:00', '2024-06-10 17:00:00', '2023-06-10 17:00:00', '2024-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (8, 2025, '2025年普通高等学校招生全国统一考试', '2025年高考', '2025-06-07 09:00:00', '2025-06-10 17:00:00', '2024-06-10 17:00:00', '2025-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (9, 2026, '2026年普通高等学校招生全国统一考试', '2026年高考', '2026-06-07 09:00:00', '2026-06-10 17:00:00', '2025-06-10 17:00:00', '2026-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (10, 2027, '2027年普通高等学校招生全国统一考试', '2027年高考', '2027-06-07 09:00:00', '2027-06-10 17:00:00', '2026-06-10 17:00:00', '2027-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (11, 2028, '2028年普通高等学校招生全国统一考试', '2028年高考', '2028-06-07 09:00:00', '2028-06-10 17:00:00', '2027-06-10 17:00:00', '2028-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (12, 2029, '2029年普通高等学校招生全国统一考试', '2029年高考', '2029-06-07 09:00:00', '2029-06-10 17:00:00', '2028-06-10 17:00:00', '2029-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (13, 2030, '2030年普通高等学校招生全国统一考试', '2030年高考', '2030-06-07 09:00:00', '2030-06-10 17:00:00', '2029-06-10 17:00:00', '2030-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (14, 2031, '2031年普通高等学校招生全国统一考试', '2031年高考', '2031-06-07 09:00:00', '2031-06-10 17:00:00', '2030-06-10 17:00:00', '2031-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (15, 2032, '2032年普通高等学校招生全国统一考试', '2032年高考', '2032-06-07 09:00:00', '2032-06-10 17:00:00', '2031-06-10 17:00:00', '2032-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (16, 2033, '2033年普通高等学校招生全国统一考试', '2033年高考', '2033-06-07 09:00:00', '2033-06-10 17:00:00', '2032-06-10 17:00:00', '2033-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (17, 2034, '2034年普通高等学校招生全国统一考试', '2034年高考', '2034-06-07 09:00:00', '2034-06-10 17:00:00', '2033-06-10 17:00:00', '2034-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (18, 2035, '2035年普通高等学校招生全国统一考试', '2035年高考', '2035-06-07 09:00:00', '2035-06-10 17:00:00', '2034-06-10 17:00:00', '2035-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (19, 2036, '2036年普通高等学校招生全国统一考试', '2036年高考', '2036-06-07 09:00:00', '2036-06-10 17:00:00', '2035-06-10 17:00:00', '2036-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (20, 2037, '2037年普通高等学校招生全国统一考试', '2037年高考', '2037-06-07 09:00:00', '2037-06-10 17:00:00', '2036-06-10 17:00:00', '2037-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (21, 2038, '2038年普通高等学校招生全国统一考试', '2038年高考', '2038-06-07 09:00:00', '2038-06-10 17:00:00', '2037-06-10 17:00:00', '2038-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (22, 2039, '2039年普通高等学校招生全国统一考试', '2039年高考', '2039-06-07 09:00:00', '2039-06-10 17:00:00', '2038-06-10 17:00:00', '2039-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (23, 2040, '2040年普通高等学校招生全国统一考试', '2040年高考', '2040-06-07 09:00:00', '2040-06-10 17:00:00', '2039-06-10 17:00:00', '2040-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (24, 2041, '2041年普通高等学校招生全国统一考试', '2041年高考', '2041-06-07 09:00:00', '2041-06-10 17:00:00', '2040-06-10 17:00:00', '2041-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (25, 2042, '2042年普通高等学校招生全国统一考试', '2042年高考', '2042-06-07 09:00:00', '2042-06-10 17:00:00', '2041-06-10 17:00:00', '2042-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (26, 2043, '2043年普通高等学校招生全国统一考试', '2043年高考', '2043-06-07 09:00:00', '2043-06-10 17:00:00', '2042-06-10 17:00:00', '2043-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (27, 2044, '2044年普通高等学校招生全国统一考试', '2044年高考', '2044-06-07 09:00:00', '2044-06-10 17:00:00', '2043-06-10 17:00:00', '2044-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (28, 2045, '2045年普通高等学校招生全国统一考试', '2045年高考', '2045-06-07 09:00:00', '2045-06-10 17:00:00', '2044-06-10 17:00:00', '2045-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (29, 2046, '2046年普通高等学校招生全国统一考试', '2046年高考', '2046-06-07 09:00:00', '2046-06-10 17:00:00', '2045-06-10 17:00:00', '2046-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (30, 2047, '2047年普通高等学校招生全国统一考试', '2047年高考', '2047-06-07 09:00:00', '2047-06-10 17:00:00', '2046-06-10 17:00:00', '2047-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (31, 2048, '2048年普通高等学校招生全国统一考试', '2048年高考', '2048-06-07 09:00:00', '2048-06-10 17:00:00', '2047-06-10 17:00:00', '2048-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (32, 2049, '2049年普通高等学校招生全国统一考试', '2049年高考', '2049-06-07 09:00:00', '2049-06-10 17:00:00', '2048-06-10 17:00:00', '2049-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (33, 2050, '2050年普通高等学校招生全国统一考试', '2050年高考', '2050-06-07 09:00:00', '2050-06-10 17:00:00', '2049-06-10 17:00:00', '2050-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (34, 2051, '2051年普通高等学校招生全国统一考试', '2051年高考', '2051-06-07 09:00:00', '2051-06-10 17:00:00', '2050-06-10 17:00:00', '2051-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (35, 2052, '2052年普通高等学校招生全国统一考试', '2052年高考', '2052-06-07 09:00:00', '2052-06-10 17:00:00', '2051-06-10 17:00:00', '2052-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (36, 2053, '2053年普通高等学校招生全国统一考试', '2053年高考', '2053-06-07 09:00:00', '2053-06-10 17:00:00', '2052-06-10 17:00:00', '2053-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (37, 2054, '2054年普通高等学校招生全国统一考试', '2054年高考', '2054-06-07 09:00:00', '2054-06-10 17:00:00', '2053-06-10 17:00:00', '2054-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (38, 2055, '2055年普通高等学校招生全国统一考试', '2055年高考', '2055-06-07 09:00:00', '2055-06-10 17:00:00', '2054-06-10 17:00:00', '2055-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (39, 2056, '2056年普通高等学校招生全国统一考试', '2056年高考', '2056-06-07 09:00:00', '2056-06-10 17:00:00', '2055-06-10 17:00:00', '2056-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (40, 2057, '2057年普通高等学校招生全国统一考试', '2057年高考', '2057-06-07 09:00:00', '2057-06-10 17:00:00', '2056-06-10 17:00:00', '2057-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (41, 2058, '2058年普通高等学校招生全国统一考试', '2058年高考', '2058-06-07 09:00:00', '2058-06-10 17:00:00', '2057-06-10 17:00:00', '2058-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (42, 2059, '2059年普通高等学校招生全国统一考试', '2059年高考', '2059-06-07 09:00:00', '2059-06-10 17:00:00', '2058-06-10 17:00:00', '2059-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (43, 2060, '2060年普通高等学校招生全国统一考试', '2060年高考', '2060-06-07 09:00:00', '2060-06-10 17:00:00', '2059-06-10 17:00:00', '2060-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (44, 2061, '2061年普通高等学校招生全国统一考试', '2061年高考', '2061-06-07 09:00:00', '2061-06-10 17:00:00', '2060-06-10 17:00:00', '2061-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (45, 2062, '2062年普通高等学校招生全国统一考试', '2062年高考', '2062-06-07 09:00:00', '2062-06-10 17:00:00', '2061-06-10 17:00:00', '2062-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (46, 2063, '2063年普通高等学校招生全国统一考试', '2063年高考', '2063-06-07 09:00:00', '2063-06-10 17:00:00', '2062-06-10 17:00:00', '2063-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (47, 2064, '2064年普通高等学校招生全国统一考试', '2064年高考', '2064-06-07 09:00:00', '2064-06-10 17:00:00', '2063-06-10 17:00:00', '2064-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (48, 2065, '2065年普通高等学校招生全国统一考试', '2065年高考', '2065-06-07 09:00:00', '2065-06-10 17:00:00', '2064-06-10 17:00:00', '2065-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (49, 2066, '2066年普通高等学校招生全国统一考试', '2066年高考', '2066-06-07 09:00:00', '2066-06-10 17:00:00', '2065-06-10 17:00:00', '2066-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (50, 2067, '2067年普通高等学校招生全国统一考试', '2067年高考', '2067-06-07 09:00:00', '2067-06-10 17:00:00', '2066-06-10 17:00:00', '2067-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (51, 2068, '2068年普通高等学校招生全国统一考试', '2068年高考', '2068-06-07 09:00:00', '2068-06-10 17:00:00', '2067-06-10 17:00:00', '2068-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (52, 2069, '2069年普通高等学校招生全国统一考试', '2069年高考', '2069-06-07 09:00:00', '2069-06-10 17:00:00', '2068-06-10 17:00:00', '2069-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (53, 2070, '2070年普通高等学校招生全国统一考试', '2070年高考', '2070-06-07 09:00:00', '2070-06-10 17:00:00', '2069-06-10 17:00:00', '2070-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (54, 2071, '2071年普通高等学校招生全国统一考试', '2071年高考', '2071-06-07 09:00:00', '2071-06-10 17:00:00', '2070-06-10 17:00:00', '2071-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (55, 2072, '2072年普通高等学校招生全国统一考试', '2072年高考', '2072-06-07 09:00:00', '2072-06-10 17:00:00', '2071-06-10 17:00:00', '2072-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (56, 2073, '2073年普通高等学校招生全国统一考试', '2073年高考', '2073-06-07 09:00:00', '2073-06-10 17:00:00', '2072-06-10 17:00:00', '2073-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (57, 2074, '2074年普通高等学校招生全国统一考试', '2074年高考', '2074-06-07 09:00:00', '2074-06-10 17:00:00', '2073-06-10 17:00:00', '2074-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (58, 2075, '2075年普通高等学校招生全国统一考试', '2075年高考', '2075-06-07 09:00:00', '2075-06-10 17:00:00', '2074-06-10 17:00:00', '2075-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (59, 2076, '2076年普通高等学校招生全国统一考试', '2076年高考', '2076-06-07 09:00:00', '2076-06-10 17:00:00', '2075-06-10 17:00:00', '2076-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (60, 2077, '2077年普通高等学校招生全国统一考试', '2077年高考', '2077-06-07 09:00:00', '2077-06-10 17:00:00', '2076-06-10 17:00:00', '2077-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (61, 2078, '2078年普通高等学校招生全国统一考试', '2078年高考', '2078-06-07 09:00:00', '2078-06-10 17:00:00', '2077-06-10 17:00:00', '2078-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (62, 2079, '2079年普通高等学校招生全国统一考试', '2079年高考', '2079-06-07 09:00:00', '2079-06-10 17:00:00', '2078-06-10 17:00:00', '2079-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (63, 2080, '2080年普通高等学校招生全国统一考试', '2080年高考', '2080-06-07 09:00:00', '2080-06-10 17:00:00', '2079-06-10 17:00:00', '2080-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (64, 2081, '2081年普通高等学校招生全国统一考试', '2081年高考', '2081-06-07 09:00:00', '2081-06-10 17:00:00', '2080-06-10 17:00:00', '2081-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (65, 2082, '2082年普通高等学校招生全国统一考试', '2082年高考', '2082-06-07 09:00:00', '2082-06-10 17:00:00', '2081-06-10 17:00:00', '2082-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (66, 2083, '2083年普通高等学校招生全国统一考试', '2083年高考', '2083-06-07 09:00:00', '2083-06-10 17:00:00', '2082-06-10 17:00:00', '2083-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (67, 2084, '2084年普通高等学校招生全国统一考试', '2084年高考', '2084-06-07 09:00:00', '2084-06-10 17:00:00', '2083-06-10 17:00:00', '2084-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (68, 2085, '2085年普通高等学校招生全国统一考试', '2085年高考', '2085-06-07 09:00:00', '2085-06-10 17:00:00', '2084-06-10 17:00:00', '2085-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (69, 2086, '2086年普通高等学校招生全国统一考试', '2086年高考', '2086-06-07 09:00:00', '2086-06-10 17:00:00', '2085-06-10 17:00:00', '2086-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (70, 2087, '2087年普通高等学校招生全国统一考试', '2087年高考', '2087-06-07 09:00:00', '2087-06-10 17:00:00', '2086-06-10 17:00:00', '2087-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (71, 2088, '2088年普通高等学校招生全国统一考试', '2088年高考', '2088-06-07 09:00:00', '2088-06-10 17:00:00', '2087-06-10 17:00:00', '2088-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (72, 2089, '2089年普通高等学校招生全国统一考试', '2089年高考', '2089-06-07 09:00:00', '2089-06-10 17:00:00', '2088-06-10 17:00:00', '2089-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (73, 2090, '2090年普通高等学校招生全国统一考试', '2090年高考', '2090-06-07 09:00:00', '2090-06-10 17:00:00', '2089-06-10 17:00:00', '2090-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (74, 2091, '2091年普通高等学校招生全国统一考试', '2091年高考', '2091-06-07 09:00:00', '2091-06-10 17:00:00', '2090-06-10 17:00:00', '2091-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (75, 2092, '2092年普通高等学校招生全国统一考试', '2092年高考', '2092-06-07 09:00:00', '2092-06-10 17:00:00', '2091-06-10 17:00:00', '2092-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (76, 2093, '2093年普通高等学校招生全国统一考试', '2093年高考', '2093-06-07 09:00:00', '2093-06-10 17:00:00', '2092-06-10 17:00:00', '2093-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (77, 2094, '2094年普通高等学校招生全国统一考试', '2094年高考', '2094-06-07 09:00:00', '2094-06-10 17:00:00', '2093-06-10 17:00:00', '2094-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (78, 2095, '2095年普通高等学校招生全国统一考试', '2095年高考', '2095-06-07 09:00:00', '2095-06-10 17:00:00', '2094-06-10 17:00:00', '2095-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (79, 2096, '2096年普通高等学校招生全国统一考试', '2096年高考', '2096-06-07 09:00:00', '2096-06-10 17:00:00', '2095-06-10 17:00:00', '2096-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (80, 2097, '2097年普通高等学校招生全国统一考试', '2097年高考', '2097-06-07 09:00:00', '2097-06-10 17:00:00', '2096-06-10 17:00:00', '2097-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (81, 2098, '2098年普通高等学校招生全国统一考试', '2098年高考', '2098-06-07 09:00:00', '2098-06-10 17:00:00', '2097-06-10 17:00:00', '2098-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (82, 2099, '2099年普通高等学校招生全国统一考试', '2099年高考', '2099-06-07 09:00:00', '2099-06-10 17:00:00', '2098-06-10 17:00:00', '2099-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (83, 2100, '2100年普通高等学校招生全国统一考试', '2100年高考', '2100-06-07 09:00:00', '2100-06-10 17:00:00', '2099-06-10 17:00:00', '2100-06-10 17:00:00', 0);
INSERT INTO `exam_date` (`id`, `exam_year`, `exam_desc`, `short_desc`, `exam_begin_date`, `exam_end_date`, `exam_year_begin_date`, `exam_year_end_date`, `is_delete`) VALUES (84, 2022, '2022年普通高等学校招生全国统一考试上海考试', '2022年上海高考', '2022-07-07 09:00:00', '2022-07-09 17:00:00', '2022-05-07 09:00:00', '2022-07-09 17:00:00', 0);
COMMIT;

-- ----------------------------
-- Table structure for send_chat
-- ----------------------------
DROP TABLE IF EXISTS `send_chat`;
CREATE TABLE `send_chat` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `chat_id` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '对话ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='发送对话';

-- ----------------------------
-- Records of send_chat
-- ----------------------------
BEGIN;
INSERT INTO `send_chat` (`id`, `chat_id`) VALUES (1, '-1001427397899');
INSERT INTO `send_chat` (`id`, `chat_id`) VALUES (2, '-1001369745974');
COMMIT;

-- ----------------------------
-- Table structure for user_template
-- ----------------------------
DROP TABLE IF EXISTS `user_template`;
CREATE TABLE `user_template` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `template_name` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '模板名称',
  `template_content` varchar(160) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '模板内容',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of user_template
-- ----------------------------
BEGIN;
INSERT INTO `user_template` (`id`, `user_id`, `template_name`, `template_content`) VALUES (1, 0, '', '现在距离{exam}还有{time}');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
