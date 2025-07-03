package com.herbertgao.telegram.bot.gaokao.business.common.util;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Date;

/**
 * DateTimeUtils 测试类
 */
public class DateTimeUtilsTest {

    @Test
    public void testFormatDurationWithMillis() {
        // 测试小于1秒的情况
        Duration shortDuration = Duration.ofMillis(500);
        assertEquals("1秒", DateTimeUtils.formatDuration(shortDuration));
        
        // 测试正好1秒
        Duration oneSecond = Duration.ofSeconds(1);
        assertEquals("1秒", DateTimeUtils.formatDuration(oneSecond));
        
        // 测试0秒
        Duration zeroDuration = Duration.ZERO;
        assertEquals("0秒", DateTimeUtils.formatDuration(zeroDuration));
        
        // 测试负数
        Duration negativeDuration = Duration.ofSeconds(-1);
        assertEquals("0秒", DateTimeUtils.formatDuration(negativeDuration));
        
        // 测试正常情况
        Duration normalDuration = Duration.ofDays(1).plusHours(2).plusMinutes(30).plusSeconds(45);
        assertEquals("1天2小时30分45秒", DateTimeUtils.formatDuration(normalDuration));
        
        // 测试只有分钟和秒
        Duration minutesSeconds = Duration.ofMinutes(30).plusSeconds(45);
        assertEquals("30分45秒", DateTimeUtils.formatDuration(minutesSeconds));
    }

    @Test
    public void testFormatDurationWithMillisPrecise() {
        // 测试精确到毫秒的方法
        Duration withMillis = Duration.ofSeconds(1).plusMillis(500);
        assertEquals("1秒500毫秒", DateTimeUtils.formatDurationWithMillis(withMillis));
        
        // 测试只有毫秒
        Duration onlyMillis = Duration.ofMillis(500);
        assertEquals("500毫秒", DateTimeUtils.formatDurationWithMillis(onlyMillis));
        
        // 测试复杂情况
        Duration complex = Duration.ofDays(1).plusHours(2).plusMinutes(30).plusSeconds(45).plusMillis(123);
        assertEquals("1天2小时30分45秒123毫秒", DateTimeUtils.formatDurationWithMillis(complex));
    }

    @Test
    public void testFormatNormal() {
        LocalDateTime dateTime = LocalDateTime.of(2023, 12, 25, 14, 30, 45);
        assertEquals("2023-12-25 14:30:45", DateTimeUtils.formatNormal(dateTime));
        
        // 测试null
        assertEquals("", DateTimeUtils.formatNormal(null));
    }

    @Test
    public void testFormat() {
        LocalDateTime dateTime = LocalDateTime.of(2023, 12, 25, 14, 30, 45);
        assertEquals("2023-12-25", DateTimeUtils.format(dateTime, "yyyy-MM-dd"));
        assertEquals("14:30", DateTimeUtils.format(dateTime, "HH:mm"));
        
        // 测试null
        assertEquals("", DateTimeUtils.format(null, "yyyy-MM-dd"));
    }

    @Test
    public void testFormatDate() {
        LocalDate date = LocalDate.of(2023, 12, 25);
        assertEquals("2023-12-25", DateTimeUtils.formatDate(date));
        
        // 测试null
        assertEquals("", DateTimeUtils.formatDate(null));
    }

    @Test
    public void testFormatTime() {
        LocalTime time = LocalTime.of(14, 30, 45);
        assertEquals("14:30:45", DateTimeUtils.formatTime(time));
        
        // 测试null
        assertEquals("", DateTimeUtils.formatTime(null));
    }

    @Test
    public void testFormatChinese() {
        LocalDateTime dateTime = LocalDateTime.of(2023, 12, 25, 14, 30, 45);
        assertEquals("2023年12月25日 14:30:45", DateTimeUtils.formatChinese(dateTime));
        
        // 测试null
        assertEquals("", DateTimeUtils.formatChinese(null));
    }

    @Test
    public void testFormatChineseDate() {
        LocalDate date = LocalDate.of(2023, 12, 25);
        assertEquals("2023年12月25日", DateTimeUtils.formatChineseDate(date));
        
        // 测试null
        assertEquals("", DateTimeUtils.formatChineseDate(null));
    }

    @Test
    public void testParseDateTime() {
        String dateTimeStr = "2023-12-25 14:30:45";
        LocalDateTime expected = LocalDateTime.of(2023, 12, 25, 14, 30, 45);
        assertEquals(expected, DateTimeUtils.parseDateTime(dateTimeStr));
        
        // 测试无效格式
        assertNull(DateTimeUtils.parseDateTime("invalid"));
        assertNull(DateTimeUtils.parseDateTime(""));
        assertNull(DateTimeUtils.parseDateTime(null));
    }

    @Test
    public void testParseDateTimeWithPattern() {
        String dateTimeStr = "2023/12/25 14:30";
        LocalDateTime expected = LocalDateTime.of(2023, 12, 25, 14, 30, 0);
        assertEquals(expected, DateTimeUtils.parseDateTime(dateTimeStr, "yyyy/MM/dd HH:mm"));
        
        // 测试无效格式
        assertNull(DateTimeUtils.parseDateTime("invalid", "yyyy/MM/dd HH:mm"));
    }

    @Test
    public void testParseDate() {
        String dateStr = "2023-12-25";
        LocalDate expected = LocalDate.of(2023, 12, 25);
        assertEquals(expected, DateTimeUtils.parseDate(dateStr));
        
        // 测试无效格式
        assertNull(DateTimeUtils.parseDate("invalid"));
        assertNull(DateTimeUtils.parseDate(""));
        assertNull(DateTimeUtils.parseDate(null));
    }

    @Test
    public void testNowAndToday() {
        LocalDateTime now = DateTimeUtils.now();
        LocalDate today = DateTimeUtils.today();
        
        assertNotNull(now);
        assertNotNull(today);
        assertEquals(LocalDate.now(), today);
    }

    @Test
    public void testCurrentTimeMillis() {
        long before = System.currentTimeMillis();
        long current = DateTimeUtils.currentTimeMillis();
        long after = System.currentTimeMillis();
        
        assertTrue(current >= before && current <= after);
    }

    @Test
    public void testTimeBetween() {
        LocalDateTime start = LocalDateTime.of(2023, 12, 25, 10, 0, 0);
        LocalDateTime end = LocalDateTime.of(2023, 12, 26, 14, 30, 45);
        
        assertEquals(1, DateTimeUtils.daysBetween(start, end));
        assertEquals(28, DateTimeUtils.hoursBetween(start, end));
        assertEquals(1710, DateTimeUtils.minutesBetween(start, end));
        assertEquals(102645, DateTimeUtils.secondsBetween(start, end));
        
        // 测试null
        assertEquals(0, DateTimeUtils.daysBetween(null, end));
        assertEquals(0, DateTimeUtils.daysBetween(start, null));
    }

    @Test
    public void testIsSameDay() {
        LocalDateTime dateTime1 = LocalDateTime.of(2023, 12, 25, 10, 0, 0);
        LocalDateTime dateTime2 = LocalDateTime.of(2023, 12, 25, 20, 0, 0);
        LocalDateTime dateTime3 = LocalDateTime.of(2023, 12, 26, 10, 0, 0);
        
        assertTrue(DateTimeUtils.isSameDay(dateTime1, dateTime2));
        assertFalse(DateTimeUtils.isSameDay(dateTime1, dateTime3));
        
        // 测试null
        assertFalse(DateTimeUtils.isSameDay(null, dateTime2));
        assertFalse(DateTimeUtils.isSameDay(dateTime1, null));
    }

    @Test
    public void testIsTodayYesterdayTomorrow() {
        LocalDateTime now = LocalDateTime.now();
        LocalDateTime yesterday = now.minusDays(1);
        LocalDateTime tomorrow = now.plusDays(1);
        
        assertTrue(DateTimeUtils.isToday(now));
        assertTrue(DateTimeUtils.isYesterday(yesterday));
        assertTrue(DateTimeUtils.isTomorrow(tomorrow));
        
        // 测试null
        assertFalse(DateTimeUtils.isToday(null));
        assertFalse(DateTimeUtils.isYesterday(null));
        assertFalse(DateTimeUtils.isTomorrow(null));
    }

    @Test
    public void testStartAndEndOfDay() {
        LocalDateTime dateTime = LocalDateTime.of(2023, 12, 25, 14, 30, 45);
        LocalDateTime startOfDay = DateTimeUtils.startOfDay(dateTime);
        LocalDateTime endOfDay = DateTimeUtils.endOfDay(dateTime);
        
        assertEquals(LocalDateTime.of(2023, 12, 25, 0, 0, 0), startOfDay);
        assertEquals(LocalDateTime.of(2023, 12, 25, 23, 59, 59, 999999999), endOfDay);
        
        // 测试null
        assertNull(DateTimeUtils.startOfDay(null));
        assertNull(DateTimeUtils.endOfDay(null));
    }

    @Test
    public void testStartAndEndOfMonth() {
        LocalDateTime dateTime = LocalDateTime.of(2023, 12, 15, 14, 30, 45);
        LocalDateTime startOfMonth = DateTimeUtils.startOfMonth(dateTime);
        LocalDateTime endOfMonth = DateTimeUtils.endOfMonth(dateTime);
        
        assertEquals(LocalDateTime.of(2023, 12, 1, 0, 0, 0), startOfMonth);
        assertEquals(LocalDateTime.of(2023, 12, 31, 23, 59, 59, 999999999), endOfMonth);
        
        // 测试null
        assertNull(DateTimeUtils.startOfMonth(null));
        assertNull(DateTimeUtils.endOfMonth(null));
    }

    @Test
    public void testStartAndEndOfYear() {
        LocalDateTime dateTime = LocalDateTime.of(2023, 6, 15, 14, 30, 45);
        LocalDateTime startOfYear = DateTimeUtils.startOfYear(dateTime);
        LocalDateTime endOfYear = DateTimeUtils.endOfYear(dateTime);
        
        assertEquals(LocalDateTime.of(2023, 1, 1, 0, 0, 0), startOfYear);
        assertEquals(LocalDateTime.of(2023, 12, 31, 23, 59, 59, 999999999), endOfYear);
        
        // 测试null
        assertNull(DateTimeUtils.startOfYear(null));
        assertNull(DateTimeUtils.endOfYear(null));
    }

    @Test
    public void testDateConversion() {
        LocalDateTime dateTime = LocalDateTime.of(2023, 12, 25, 14, 30, 45);
        Date date = DateTimeUtils.toDate(dateTime);
        LocalDateTime converted = DateTimeUtils.fromDate(date);
        
        assertEquals(dateTime, converted);
        
        // 测试null
        assertNull(DateTimeUtils.toDate(null));
        assertNull(DateTimeUtils.fromDate(null));
    }
} 