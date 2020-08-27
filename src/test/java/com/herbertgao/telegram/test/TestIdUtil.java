package com.herbertgao.telegram.test;


import com.herbertgao.telegram.util.IdUtil;
import org.junit.Test;

public class TestIdUtil {

    @Test
    public void testIdSnowflake() {
        System.out.println(IdUtil.getId());
    }
}
