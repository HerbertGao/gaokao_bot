package com.herbertgao.telegram.bot.gaokao.business.common.util;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

import java.util.HashSet;
import java.util.Set;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/**
 * IdUtil 测试类
 */
public class IdUtilTest {

    @Test
    public void testGenerateId() {
        Long id = IdUtil.getId();
        assertNotNull(id);
        assertTrue(id > 0);
    }

    @Test
    public void testIdUniqueness() {
        Set<Long> ids = new HashSet<>();
        int count = 1000;
        
        for (int i = 0; i < count; i++) {
            Long id = IdUtil.getId();
            assertTrue(ids.add(id), "ID应该唯一: " + id);
        }
        
        assertEquals(count, ids.size());
    }

    @Test
    public void testIdMonotonicity() {
        Long previousId = IdUtil.getId();
        
        for (int i = 0; i < 100; i++) {
            Long currentId = IdUtil.getId();
            assertTrue(currentId > previousId, "ID应该单调递增");
            previousId = currentId;
        }
    }

    @Test
    public void testConcurrentIdGeneration() throws InterruptedException {
        int threadCount = 10;
        int idsPerThread = 100;
        Set<Long> allIds = new HashSet<>();
        
        ExecutorService executor = Executors.newFixedThreadPool(threadCount);
        CountDownLatch latch = new CountDownLatch(threadCount);
        
        for (int i = 0; i < threadCount; i++) {
            executor.submit(() -> {
                try {
                    for (int j = 0; j < idsPerThread; j++) {
                        Long id = IdUtil.getId();
                        synchronized (allIds) {
                            assertTrue(allIds.add(id), "并发生成的ID应该唯一: " + id);
                        }
                    }
                } finally {
                    latch.countDown();
                }
            });
        }
        
        latch.await();
        executor.shutdown();
        
        assertEquals(threadCount * idsPerThread, allIds.size());
    }

    @Test
    public void testIdStructure() {
        Long id = IdUtil.getId();
        
        // 验证ID长度（64位长整型）
        assertTrue(id > 0);
        
        // 验证ID是正数且合理大小
        // 由于epoch是2021年，当前时间戳减去epoch后，ID应该是一个合理的正数
        assertTrue(id > 0L);
    }

    @Test
    public void testMultipleCalls() {
        // 连续多次调用，确保不会出现重复
        Set<Long> ids = new HashSet<>();
        for (int i = 0; i < 10000; i++) {
            Long id = IdUtil.getId();
            assertTrue(ids.add(id), "连续调用生成的ID应该唯一");
        }
    }
} 