package com.herbertgao.telegram.bot.gaokao;

import com.herbertgao.telegram.bot.gaokao.business.service.GaokaoBotService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.inlinequery.InlineQuery;

/**
 * gaokao_bot
 *
 * @author HerbertGao
 * @date 2023-06-07
 */
@Slf4j
@Component
public class GaokaoBot extends TelegramLongPollingBot {

    private final String botUsername;

    @Autowired
    private GaokaoBotService gaokaoBotService;

    public GaokaoBot(@Value("${telegram.bot.username}") String botUsername,
                     @Value("${telegram.bot.token}") String botToken) {
        super(botToken);
        this.botUsername = botUsername;
    }


    @Override
    public void onUpdateReceived(Update update) {
        log.debug(update.toString());

        if (update.hasInlineQuery()) {
            InlineQuery inlineQuery = update.getInlineQuery();
            gaokaoBotService.inlineQuery(inlineQuery);
        } else if (update.hasMessage()) {
            Message message = update.getMessage();
            gaokaoBotService.message(message);
        }
    }

    @Override
    public String getBotUsername() {
        return botUsername;
    }
}
