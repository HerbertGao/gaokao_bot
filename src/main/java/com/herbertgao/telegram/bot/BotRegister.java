package com.herbertgao.telegram.bot;

import com.herbertgao.telegram.bot.gaokao.GaokaoBot;
import jakarta.annotation.PostConstruct;
import lombok.SneakyThrows;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.meta.TelegramBotsApi;
import org.telegram.telegrambots.updatesreceivers.DefaultBotSession;

@Component
public class BotRegister {

    @Autowired
    GaokaoBot gaokaoBot;

    @SneakyThrows
    @PostConstruct
    public void init() {
        TelegramBotsApi telegramBotsApi = new TelegramBotsApi(DefaultBotSession.class);
        telegramBotsApi.registerBot(gaokaoBot);
    }

}
