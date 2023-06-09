package com.herbertgao.telegram.bot;

import com.herbertgao.telegram.bot.gaokao.GaokaoBot;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.meta.TelegramBotsApi;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;
import org.telegram.telegrambots.updatesreceivers.DefaultBotSession;

@Slf4j
@Component
public class BotRegister {

    @Autowired
    GaokaoBot gaokaoBot;

    @PostConstruct
    public void init() {
        try {
            TelegramBotsApi telegramBotsApi = new TelegramBotsApi(DefaultBotSession.class);
            telegramBotsApi.registerBot(gaokaoBot);
        } catch (TelegramApiException e) {
            log.error(e.getMessage());
        }
    }

}
