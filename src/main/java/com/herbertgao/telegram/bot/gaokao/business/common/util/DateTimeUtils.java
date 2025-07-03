package com.herbertgao.telegram.bot.gaokao.business.common.util;

import java.time.*;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeParseException;
import java.time.temporal.ChronoUnit;
import java.util.Date;

/**
 * 日期时间工具类
 * 替代 Hutool 的 DateUtil 相关功能，提供统一的时间处理工具
 *
 * @author HerbertGao
 */
public class DateTimeUtils {

    // 常用日期时间格式
    public static final String PATTERN_NORMAL = "yyyy-MM-dd HH:mm:ss";
    public static final String PATTERN_DATE = "yyyy-MM-dd";
    public static final String PATTERN_TIME = "HH:mm:ss";
    public static final String PATTERN_DATETIME = "yyyy-MM-dd HH:mm";
    public static final String PATTERN_CN_DATE = "yyyy年MM月dd日";
    public static final String PATTERN_CN_DATETIME = "yyyy年MM月dd日 HH:mm:ss";
    public static final String PATTERN_ISO = "yyyy-MM-dd'T'HH:mm:ss";

    // 预编译的格式化器，提高性能
    private static final DateTimeFormatter NORMAL_FORMATTER = DateTimeFormatter.ofPattern(PATTERN_NORMAL);
    private static final DateTimeFormatter DATE_FORMATTER = DateTimeFormatter.ofPattern(PATTERN_DATE);
    private static final DateTimeFormatter TIME_FORMATTER = DateTimeFormatter.ofPattern(PATTERN_TIME);
    private static final DateTimeFormatter DATETIME_FORMATTER = DateTimeFormatter.ofPattern(PATTERN_DATETIME);
    private static final DateTimeFormatter CN_DATE_FORMATTER = DateTimeFormatter.ofPattern(PATTERN_CN_DATE);
    private static final DateTimeFormatter CN_DATETIME_FORMATTER = DateTimeFormatter.ofPattern(PATTERN_CN_DATETIME);
    private static final DateTimeFormatter ISO_FORMATTER = DateTimeFormatter.ofPattern(PATTERN_ISO);

    /**
     * 格式化LocalDateTime为yyyy-MM-dd HH:mm:ss字符串
     *
     * @param dateTime LocalDateTime对象
     * @return 格式化后的字符串，若为null返回空串
     */
    public static String formatNormal(LocalDateTime dateTime) {
        if (dateTime == null) {
            return "";
        }
        return dateTime.format(NORMAL_FORMATTER);
    }

    /**
     * 格式化LocalDateTime为指定格式字符串
     *
     * @param dateTime LocalDateTime对象
     * @param pattern  格式模式
     * @return 格式化后的字符串，若为null返回空串
     */
    public static String format(LocalDateTime dateTime, String pattern) {
        if (dateTime == null) {
            return "";
        }
        return dateTime.format(DateTimeFormatter.ofPattern(pattern));
    }

    /**
     * 格式化LocalDate为yyyy-MM-dd字符串
     *
     * @param date LocalDate对象
     * @return 格式化后的字符串，若为null返回空串
     */
    public static String formatDate(LocalDate date) {
        if (date == null) {
            return "";
        }
        return date.format(DATE_FORMATTER);
    }

    /**
     * 格式化LocalTime为HH:mm:ss字符串
     *
     * @param time LocalTime对象
     * @return 格式化后的字符串，若为null返回空串
     */
    public static String formatTime(LocalTime time) {
        if (time == null) {
            return "";
        }
        return time.format(TIME_FORMATTER);
    }

    /**
     * 格式化LocalDateTime为中文日期时间格式
     *
     * @param dateTime LocalDateTime对象
     * @return 格式化后的字符串，若为null返回空串
     */
    public static String formatChinese(LocalDateTime dateTime) {
        if (dateTime == null) {
            return "";
        }
        return dateTime.format(CN_DATETIME_FORMATTER);
    }

    /**
     * 格式化LocalDate为中文日期格式
     *
     * @param date LocalDate对象
     * @return 格式化后的字符串，若为null返回空串
     */
    public static String formatChineseDate(LocalDate date) {
        if (date == null) {
            return "";
        }
        return date.format(CN_DATE_FORMATTER);
    }

    /**
     * 解析字符串为LocalDateTime
     *
     * @param dateTimeStr 日期时间字符串
     * @return LocalDateTime对象，解析失败返回null
     */
    public static LocalDateTime parseDateTime(String dateTimeStr) {
        if (dateTimeStr == null || dateTimeStr.trim().isEmpty()) {
            return null;
        }
        try {
            return LocalDateTime.parse(dateTimeStr.trim(), NORMAL_FORMATTER);
        } catch (DateTimeParseException e) {
            return null;
        }
    }

    /**
     * 解析字符串为LocalDateTime（指定格式）
     *
     * @param dateTimeStr 日期时间字符串
     * @param pattern     格式模式
     * @return LocalDateTime对象，解析失败返回null
     */
    public static LocalDateTime parseDateTime(String dateTimeStr, String pattern) {
        if (dateTimeStr == null || dateTimeStr.trim().isEmpty()) {
            return null;
        }
        try {
            return LocalDateTime.parse(dateTimeStr.trim(), DateTimeFormatter.ofPattern(pattern));
        } catch (DateTimeParseException e) {
            return null;
        }
    }

    /**
     * 解析字符串为LocalDate
     *
     * @param dateStr 日期字符串
     * @return LocalDate对象，解析失败返回null
     */
    public static LocalDate parseDate(String dateStr) {
        if (dateStr == null || dateStr.trim().isEmpty()) {
            return null;
        }
        try {
            return LocalDate.parse(dateStr.trim(), DATE_FORMATTER);
        } catch (DateTimeParseException e) {
            return null;
        }
    }

    /**
     * 获取当前时间
     *
     * @return 当前LocalDateTime
     */
    public static LocalDateTime now() {
        return LocalDateTime.now();
    }

    /**
     * 获取当前日期
     *
     * @return 当前LocalDate
     */
    public static LocalDate today() {
        return LocalDate.now();
    }

    /**
     * 获取当前时间戳（毫秒）
     *
     * @return 当前时间戳
     */
    public static long currentTimeMillis() {
        return System.currentTimeMillis();
    }

    /**
     * 获取两个时间之间的天数差
     *
     * @param start 开始时间
     * @param end   结束时间
     * @return 天数差
     */
    public static long daysBetween(LocalDateTime start, LocalDateTime end) {
        if (start == null || end == null) {
            return 0;
        }
        return ChronoUnit.DAYS.between(start, end);
    }

    /**
     * 获取两个时间之间的小时数差
     *
     * @param start 开始时间
     * @param end   结束时间
     * @return 小时数差
     */
    public static long hoursBetween(LocalDateTime start, LocalDateTime end) {
        if (start == null || end == null) {
            return 0;
        }
        return ChronoUnit.HOURS.between(start, end);
    }

    /**
     * 获取两个时间之间的分钟数差
     *
     * @param start 开始时间
     * @param end   结束时间
     * @return 分钟数差
     */
    public static long minutesBetween(LocalDateTime start, LocalDateTime end) {
        if (start == null || end == null) {
            return 0;
        }
        return ChronoUnit.MINUTES.between(start, end);
    }

    /**
     * 获取两个时间之间的秒数差
     *
     * @param start 开始时间
     * @param end   结束时间
     * @return 秒数差
     */
    public static long secondsBetween(LocalDateTime start, LocalDateTime end) {
        if (start == null || end == null) {
            return 0;
        }
        return ChronoUnit.SECONDS.between(start, end);
    }

    /**
     * 判断是否为同一天
     *
     * @param dateTime1 时间1
     * @param dateTime2 时间2
     * @return 是否为同一天
     */
    public static boolean isSameDay(LocalDateTime dateTime1, LocalDateTime dateTime2) {
        if (dateTime1 == null || dateTime2 == null) {
            return false;
        }
        return dateTime1.toLocalDate().equals(dateTime2.toLocalDate());
    }

    /**
     * 判断是否为今天
     *
     * @param dateTime 时间
     * @return 是否为今天
     */
    public static boolean isToday(LocalDateTime dateTime) {
        if (dateTime == null) {
            return false;
        }
        return isSameDay(dateTime, now());
    }

    /**
     * 判断是否为昨天
     *
     * @param dateTime 时间
     * @return 是否为昨天
     */
    public static boolean isYesterday(LocalDateTime dateTime) {
        if (dateTime == null) {
            return false;
        }
        return isSameDay(dateTime, now().minusDays(1));
    }

    /**
     * 判断是否为明天
     *
     * @param dateTime 时间
     * @return 是否为明天
     */
    public static boolean isTomorrow(LocalDateTime dateTime) {
        if (dateTime == null) {
            return false;
        }
        return isSameDay(dateTime, now().plusDays(1));
    }

    /**
     * 获取时间开始（00:00:00）
     *
     * @param dateTime 时间
     * @return 当天开始时间
     */
    public static LocalDateTime startOfDay(LocalDateTime dateTime) {
        if (dateTime == null) {
            return null;
        }
        return dateTime.toLocalDate().atStartOfDay();
    }

    /**
     * 获取时间结束（23:59:59.999999999）
     *
     * @param dateTime 时间
     * @return 当天结束时间
     */
    public static LocalDateTime endOfDay(LocalDateTime dateTime) {
        if (dateTime == null) {
            return null;
        }
        return dateTime.toLocalDate().atTime(23, 59, 59, 999999999);
    }

    /**
     * 获取月份开始时间
     *
     * @param dateTime 时间
     * @return 月份开始时间
     */
    public static LocalDateTime startOfMonth(LocalDateTime dateTime) {
        if (dateTime == null) {
            return null;
        }
        return dateTime.toLocalDate().withDayOfMonth(1).atStartOfDay();
    }

    /**
     * 获取月份结束时间
     *
     * @param dateTime 时间
     * @return 月份结束时间
     */
    public static LocalDateTime endOfMonth(LocalDateTime dateTime) {
        if (dateTime == null) {
            return null;
        }
        return dateTime.toLocalDate().withDayOfMonth(
                dateTime.toLocalDate().lengthOfMonth()
        ).atTime(23, 59, 59, 999999999);
    }

    /**
     * 获取年份开始时间
     *
     * @param dateTime 时间
     * @return 年份开始时间
     */
    public static LocalDateTime startOfYear(LocalDateTime dateTime) {
        if (dateTime == null) {
            return null;
        }
        return dateTime.toLocalDate().withDayOfYear(1).atStartOfDay();
    }

    /**
     * 获取年份结束时间
     *
     * @param dateTime 时间
     * @return 年份结束时间
     */
    public static LocalDateTime endOfYear(LocalDateTime dateTime) {
        if (dateTime == null) {
            return null;
        }
        return dateTime.toLocalDate().withDayOfYear(
                dateTime.toLocalDate().lengthOfYear()
        ).atTime(23, 59, 59, 999999999);
    }

    /**
     * Date转LocalDateTime
     *
     * @param date Date对象
     * @return LocalDateTime对象
     */
    public static LocalDateTime fromDate(Date date) {
        if (date == null) {
            return null;
        }
        return date.toInstant().atZone(ZoneId.systemDefault()).toLocalDateTime();
    }

    /**
     * LocalDateTime转Date
     *
     * @param dateTime LocalDateTime对象
     * @return Date对象
     */
    public static Date toDate(LocalDateTime dateTime) {
        if (dateTime == null) {
            return null;
        }
        return Date.from(dateTime.atZone(ZoneId.systemDefault()).toInstant());
    }

    /**
     * 格式化时间间隔为中文格式
     * 格式：天时分秒，如果某项为0则不显示
     * 例如：350天23小时59分59秒、18天3分
     * <p>
     * 注意：如果总时长小于1秒但大于0，会显示为"1秒"
     *
     * @param duration 时间间隔
     * @return 格式化后的字符串
     */
    public static String formatDuration(Duration duration) {
        if (duration.isNegative()) {
            return "0秒";
        }

        long totalSeconds = duration.getSeconds();
        long nanoSeconds = duration.getNano();
        
        // 如果有毫秒部分且总时长小于1秒，则显示为1秒
        if (totalSeconds == 0 && nanoSeconds > 0) {
            return "1秒";
        }
        
        // 如果总时长为0，返回0秒
        if (totalSeconds == 0 && nanoSeconds == 0) {
            return "0秒";
        }

        long days = totalSeconds / 86400; // 24 * 60 * 60
        long hours = (totalSeconds % 86400) / 3600;
        long minutes = (totalSeconds % 3600) / 60;
        long seconds = totalSeconds % 60;

        StringBuilder result = new StringBuilder();
        
        if (days > 0) {
            result.append(days).append("天");
        }
        if (hours > 0) {
            result.append(hours).append("小时");
        }
        if (minutes > 0) {
            result.append(minutes).append("分");
        }
        if (seconds > 0 || result.isEmpty()) {
            result.append(seconds).append("秒");
        }

        return result.toString();
    }

    /**
     * 格式化时间间隔为中文格式（精确到毫秒）
     * 格式：天时分秒毫秒，如果某项为0则不显示
     * 例如：350天23小时59分59秒500毫秒、18天3分200毫秒
     *
     * @param duration 时间间隔
     * @return 格式化后的字符串
     */
    public static String formatDurationWithMillis(Duration duration) {
        if (duration.isNegative()) {
            return "0秒";
        }

        long totalSeconds = duration.getSeconds();
        long nanoSeconds = duration.getNano();
        long millis = nanoSeconds / 1_000_000; // 转换为毫秒
        
        // 如果总时长为0，返回0秒
        if (totalSeconds == 0 && millis == 0) {
            return "0秒";
        }

        long days = totalSeconds / 86400; // 24 * 60 * 60
        long hours = (totalSeconds % 86400) / 3600;
        long minutes = (totalSeconds % 3600) / 60;
        long seconds = totalSeconds % 60;

        StringBuilder result = new StringBuilder();
        
        if (days > 0) {
            result.append(days).append("天");
        }
        if (hours > 0) {
            result.append(hours).append("小时");
        }
        if (minutes > 0) {
            result.append(minutes).append("分");
        }
        if (seconds > 0) {
            result.append(seconds).append("秒");
        }
        if (millis > 0) {
            result.append(millis).append("毫秒");
        }
        
        // 如果没有任何时间单位，说明只有毫秒
        if (result.isEmpty()) {
            result.append(millis).append("毫秒");
        }

        return result.toString();
    }
} 