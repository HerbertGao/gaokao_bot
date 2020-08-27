package com.herbertgao.telegram.business.service;

import com.herbertgao.telegram.bot.Command;
import com.herbertgao.telegram.bot.Config;
import com.herbertgao.telegram.util.TelegramBotUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.telegram.telegrambots.meta.api.methods.AnswerInlineQuery;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.meta.api.objects.User;
import org.telegram.telegrambots.meta.api.objects.inlinequery.InlineQuery;
import org.telegram.telegrambots.meta.bots.AbsSender;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

/**
 * @program: gaokao_bot
 * @description: TelegramBot Service Implements
 * @author: HerbertGao
 * @create: 2019-06-08 23:55
 **/
@Service
public class MyBotService {

    private final static Logger logger = LoggerFactory.getLogger(MyBotService.class);

    @Autowired
    private AbsSender absSender;

    @Autowired
    private InlineQueryService inlineQueryService;
    @Autowired
    private MessageService messageService;

    /**
     * Inline Query
     *
     * @param inlineQuery
     */
    public void inlineQuery(InlineQuery inlineQuery) {
        logger.debug(inlineQuery.toString());

        try {
            AnswerInlineQuery results = inlineQueryService.answerInlineQuery(inlineQuery);
            absSender.execute(results);
        } catch (TelegramApiException e) {
            e.printStackTrace();
        }
    }

    /**
     * Message
     *
     * @param message
     */
    public void message(Message message) {
        logger.debug(message.toString());

        Long chatId = message.getChatId();
        User user = message.getFrom();

        if (message.hasText()) {
            String text = TelegramBotUtil.getTextByMessage(message, null);

            if (message.isCommand()) {

                if (text.startsWith(Command.COUNTDOWN_COMMAND)) {
                    String msg = messageService.getCountDownCommandMessage();

                    try {
                        absSender.execute(new SendMessage(chatId, msg).enableHtml(true).setReplyToMessageId(message.getMessageId()));
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (text.startsWith(Command.LIST_COMMAND)) {
                    String msg = messageService.getListCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg).enableHtml(true).enableHtml(true).setReplyToMessageId(message.getMessageId()));
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (text.startsWith(Command.ADD_COMMAND)) {
                    String msg = messageService.getAddCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg).enableHtml(true).setReplyToMessageId(message.getMessageId()));
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (text.startsWith(Command.REMOVE_COMMAND)) {
                    String msg = messageService.getRemoveCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg).enableHtml(true).setReplyToMessageId(message.getMessageId()));
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (text.startsWith(Command.CUSTOMIZE_COMMAND)) {
                    String msg = messageService.getCustomizeCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg).enableHtml(true).setReplyToMessageId(message.getMessageId()));
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (text.startsWith(Command.RENAME_COMMAND)) {
                    String msg = messageService.getRenameCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg).enableHtml(true).setReplyToMessageId(message.getMessageId()));
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                }

            }
        }
    }
}
