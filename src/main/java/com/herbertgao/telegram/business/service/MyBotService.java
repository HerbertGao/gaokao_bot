package com.herbertgao.telegram.business.service;

import com.herbertgao.telegram.bot.Command;
import com.herbertgao.telegram.util.TelegramBotUtil;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Lazy;
import org.springframework.stereotype.Service;
import org.telegram.telegrambots.meta.api.methods.AnswerInlineQuery;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.meta.api.objects.User;
import org.telegram.telegrambots.meta.api.objects.inlinequery.InlineQuery;
import org.telegram.telegrambots.meta.bots.AbsSender;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

/**
 * 我的机器人服务
 *
 * @author HerbertGao
 * @date 2019-06-08
 */
@Slf4j
@Service
public class MyBotService {

    @Lazy
    @Autowired
    private AbsSender absSender;

    @Autowired
    private InlineQueryService inlineQueryService;
    @Autowired
    private MessageService messageService;

    /**
     * Inline Query
     *
     * @param inlineQuery 内联查询
     */
    public void inlineQuery(InlineQuery inlineQuery) {
        log.debug(inlineQuery.toString());

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
     * @param message 消息
     */
    public void message(Message message) {
        log.debug(message.toString());

        String chatId = message.getChatId().toString();
        User user = message.getFrom();

        if (message.hasText()) {
            String text = TelegramBotUtil.getTextByMessage(message, null);

            if (message.isCommand()) {

                if (TelegramBotUtil.isMatchCommand(text, Command.COUNTDOWN_COMMAND)) {
                    String msg = messageService.getCountDownCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg) {
                            {
                                enableHtml(true);
                                setReplyToMessageId(message.getMessageId());
                            }
                        });
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (TelegramBotUtil.isMatchCommand(text, Command.LIST_COMMAND)) {
                    String msg = messageService.getListCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg) {
                            {
                                enableHtml(true);
                                setReplyToMessageId(message.getMessageId());
                            }
                        });
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (TelegramBotUtil.isMatchCommand(text, Command.ADD_COMMAND)) {
                    String msg = messageService.getAddCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg) {
                            {
                                enableHtml(true);
                                setReplyToMessageId(message.getMessageId());
                            }
                        });
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (TelegramBotUtil.isMatchCommand(text, Command.REMOVE_COMMAND)) {
                    String msg = messageService.getRemoveCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg) {
                            {
                                enableHtml(true);
                                setReplyToMessageId(message.getMessageId());
                            }
                        });
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (TelegramBotUtil.isMatchCommand(text, Command.CUSTOMIZE_COMMAND)) {
                    String msg = messageService.getCustomizeCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg) {
                            {
                                enableHtml(true);
                                setReplyToMessageId(message.getMessageId());
                            }
                        });
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                } else if (TelegramBotUtil.isMatchCommand(text, Command.RENAME_COMMAND)) {
                    String msg = messageService.getRenameCommandMessage(message);

                    try {
                        absSender.execute(new SendMessage(chatId, msg) {
                            {
                                enableHtml(true);
                                setReplyToMessageId(message.getMessageId());
                            }
                        });
                    } catch (TelegramApiException e) {
                        e.printStackTrace();
                    }
                }

            }
        }
    }
}
