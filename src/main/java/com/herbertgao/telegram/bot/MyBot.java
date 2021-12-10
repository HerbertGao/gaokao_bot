package com.herbertgao.telegram.bot;

import com.herbertgao.telegram.business.service.MyBotService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.inlinequery.InlineQuery;

/**
 * 我的Bot
 *
 * @author HerbertGao
 * @date 2019-06-08
 */
@Slf4j
@Component
public class MyBot extends TelegramLongPollingBot {

    @Autowired
    private MyBotService myBotService;

    @Override
    public void onUpdateReceived(Update update) {
        log.debug(update.toString());

        if (update.hasInlineQuery()) {
            InlineQuery inlineQuery = update.getInlineQuery();
            myBotService.inlineQuery(inlineQuery);
        } else if (update.hasMessage()) {
            Message message = update.getMessage();
            myBotService.message(message);
        }
    }

    @Override
    public String getBotUsername() {
        return Config.getUsername();
    }

    @Override
    public String getBotToken() {
        return Config.getToken();
    }

}
