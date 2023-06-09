package com.herbertgao.telegram.bot.gaokao.business.common.constant;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

@Configuration
public class Config {

    public static String botUsername;

    @Value("${telegram.bot.username}")
    public void setBotUsername(String botUsername) {
        Config.botUsername = botUsername;
    }

}
