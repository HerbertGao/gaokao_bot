package com.herbertgao.telegram.bot.gaokao.business.common.util;

import java.util.concurrent.TimeUnit;

/**
 * 线程工具类
 * 用于替换Hutool的ThreadUtil功能
 *
 * @author HerbertGao
 * @date 2025-07-02
 */
public class ThreadUtils {

    /**
     * 安全休眠
     * 替换Hutool的ThreadUtil.safeSleep方法
     * 在休眠过程中如果被中断，会正确处理中断状态
     *
     * @param milliseconds 休眠时间（毫秒）
     */
    public static void safeSleep(long milliseconds) {
        try {
            TimeUnit.MILLISECONDS.sleep(milliseconds);
        } catch (InterruptedException e) {
            // 恢复中断状态，这是Java并发编程的最佳实践
            Thread.currentThread().interrupt();
        }
    }

    /**
     * 安全休眠（秒）
     *
     * @param seconds 休眠时间（秒）
     */
    public static void safeSleepSeconds(long seconds) {
        safeSleep(seconds * 1000);
    }

    /**
     * 安全休眠（分钟）
     *
     * @param minutes 休眠时间（分钟）
     */
    public static void safeSleepMinutes(long minutes) {
        safeSleep(minutes * 60 * 1000);
    }

    /**
     * 获取当前线程名称
     *
     * @return 当前线程名称
     */
    public static String getCurrentThreadName() {
        return Thread.currentThread().getName();
    }

    /**
     * 检查当前线程是否被中断
     *
     * @return 是否被中断
     */
    public static boolean isInterrupted() {
        return Thread.currentThread().isInterrupted();
    }

    /**
     * 让出CPU时间片
     * 让其他线程有机会执行
     */
    public static void yield() {
        Thread.yield();
    }
} 