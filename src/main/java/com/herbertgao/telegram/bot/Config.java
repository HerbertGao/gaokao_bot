package com.herbertgao.telegram.bot;

import org.springframework.boot.context.properties.ConfigurationProperties;

/**
 * @program: gaokao_bot
 * @description: TelegramConfig
 * @author: HerbertGao
 * @create: 2019-06-08 22:51
 **/
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
