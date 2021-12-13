package com.herbertgao.telegram.util;

import cn.hutool.core.lang.Snowflake;

/**
 * IDUtil
 *
 * @author HerbertGao
 * @date 2021-12-10
 */
public class IdUtil {

    static Snowflake snowflake = cn.hutool.core.util.IdUtil.getSnowflake(1, 1);

    public static Long getId() {
        return snowflake.nextId();
    }

}
