package com.herbertgao.telegram.util;

import cn.hutool.core.lang.Snowflake;

public class IdUtil {

    static Snowflake snowflake = cn.hutool.core.util.IdUtil.createSnowflake(1, 1);

    public static Long getId() {
        return snowflake.nextId();
    }

}
