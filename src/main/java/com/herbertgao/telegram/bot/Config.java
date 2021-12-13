package com.herbertgao.telegram.bot;

import org.springframework.boot.context.properties.ConfigurationProperties;

/**
 * 配置
 *
 * @author HerbertGao
 * @date 2019-06-08
 */
@ConfigurationProperties(prefix = "telegram.bot")
public class Config {

    private static String username;
    private static String token;

    public static String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        Config.username = username;
    }

    public static String getToken() {
        return token;
    }

    public void setToken(String token) {
        Config.token = token;
    }
}
