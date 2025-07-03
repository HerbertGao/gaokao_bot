package com.herbertgao.telegram.bot;

import com.herbertgao.telegram.bot.gaokao.GaokaoBot;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.longpolling.TelegramBotsLongPollingApplication;

@Slf4j
@Component
public class BotRegister {

    @Autowired
    GaokaoBot gaokaoBot;

    @PostConstruct
    public void init() {
        try {
            TelegramBotsLongPollingApplication botsApplication = new TelegramBotsLongPollingApplication();
            botsApplication.registerBot(gaokaoBot.getBotToken(), gaokaoBot);
        } catch (Exception e) {
            log.error(e.getMessage());
        }
    }

}
