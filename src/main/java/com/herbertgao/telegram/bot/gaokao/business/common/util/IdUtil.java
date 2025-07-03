package com.herbertgao.telegram.bot.gaokao.business.common.util;

import java.time.Instant;
import java.util.concurrent.atomic.AtomicLong;

/**
 * ID Util
 * 替代Hutool的Snowflake ID生成器
 *
 * @author HerbertGao
 * @date 2021-12-14
 */
public class IdUtil {

    private static final long EPOCH = 1288834974657L; // 与Hutool保持一致
    private static final long WORKER_ID_BITS = 5L;
    private static final long DATACENTER_ID_BITS = 5L;
    private static final long SEQUENCE_BITS = 12L;
    
    private static final long MAX_WORKER_ID = ~(-1L << WORKER_ID_BITS);
    private static final long MAX_DATACENTER_ID = ~(-1L << DATACENTER_ID_BITS);
    
    private static final long WORKER_ID_SHIFT = SEQUENCE_BITS;
    private static final long DATACENTER_ID_SHIFT = SEQUENCE_BITS + WORKER_ID_BITS;
    private static final long TIMESTAMP_LEFT_SHIFT = SEQUENCE_BITS + WORKER_ID_BITS + DATACENTER_ID_BITS;
    
    private static final long SEQUENCE_MASK = ~(-1L << SEQUENCE_BITS);
    
    private static final long WORKER_ID = 1L;
    private static final long DATACENTER_ID = 1L;
    
    private static final AtomicLong sequence = new AtomicLong(0L);
    private static long lastTimestamp = -1L;

    public static Long getId() {
        return nextId();
    }
    
    private static synchronized long nextId() {
        long timestamp = timeGen();
        
        if (timestamp < lastTimestamp) {
            // 时钟回拨，等待直到时间追上
            timestamp = tilNextMillis(lastTimestamp);
        }
        
        if (lastTimestamp == timestamp) {
            long currentSequence = sequence.incrementAndGet() & SEQUENCE_MASK;
            if (currentSequence == 0) {
                timestamp = tilNextMillis(lastTimestamp);
            }
        } else {
            sequence.set(0L);
        }
        
        lastTimestamp = timestamp;
        
        return ((timestamp - EPOCH) << TIMESTAMP_LEFT_SHIFT) |
                (DATACENTER_ID << DATACENTER_ID_SHIFT) |
                (WORKER_ID << WORKER_ID_SHIFT) |
                sequence.get();
    }

    private static long tilNextMillis(long lastTimestamp) {
        long timestamp = timeGen();
        while (timestamp <= lastTimestamp) {
            timestamp = timeGen();
        }
        return timestamp;
    }
    
    private static long timeGen() {
        return Instant.now().toEpochMilli();
    }
}
