package com.herbertgao.telegram.bot.gaokao;

import com.herbertgao.telegram.bot.gaokao.business.service.GaokaoBotService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.longpolling.interfaces.LongPollingUpdateConsumer;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.inlinequery.InlineQuery;
import org.telegram.telegrambots.meta.api.objects.message.Message;

import java.util.List;

/**
 * gaokao_bot
 *
 * @author HerbertGao
 * @date 2023-06-07
 */
@Slf4j
@Component
public class GaokaoBot implements LongPollingUpdateConsumer {

    private final String botUsername;
    private final String botToken;

    @Autowired
    private GaokaoBotService gaokaoBotService;

    public GaokaoBot(@Value("${telegram.bot.username}") String botUsername,
                     @Value("${telegram.bot.token}") String botToken) {
        this.botUsername = botUsername;
        this.botToken = botToken;
    }

    @Override
    public void consume(List<Update> updates) {
        for (Update update : updates) {
            log.debug(update.toString());

            if (update.hasInlineQuery()) {
                InlineQuery inlineQuery = update.getInlineQuery();
                gaokaoBotService.inlineQuery(inlineQuery);
            } else if (update.hasMessage()) {
                Message message = update.getMessage();
                gaokaoBotService.message(message);
            }
        }
    }

    public String getBotUsername() {
        return botUsername;
    }

    public String getBotToken() {
        return botToken;
    }
}
