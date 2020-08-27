package com.herbertgao.telegram.bot;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.ComponentScan;
import org.telegram.telegrambots.ApiContextInitializer;
import org.telegram.telegrambots.meta.TelegramBotsApi;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

/**
 * @program: gaokao_bot
 * @description: Main Class
 * @author: HerbertGao
 * @create: 2019-06-08 21:49
 **/
@SpringBootApplication
@EnableConfigurationProperties(Config.class)
@ComponentScan(basePackages = {"com.herbertgao.telegram"})
@MapperScan(basePackages = "com.herbertgao.telegram.database.mapper")
public class Main implements ApplicationRunner {

    @Autowired
    private MyBot bot;

    public static void main(String[] args) {
        ApiContextInitializer.init();
        SpringApplication.run(Main.class, args);
    }

    @Override
    public void run(ApplicationArguments args) {
        TelegramBotsApi botsApi = new TelegramBotsApi();

        try {
            botsApi.registerBot(bot);
        } catch (TelegramApiException e) {
            e.printStackTrace();
        }
    }
}
