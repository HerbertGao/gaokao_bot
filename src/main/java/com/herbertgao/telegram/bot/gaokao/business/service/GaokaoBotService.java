package com.herbertgao.telegram.bot.gaokao.business.service;

import com.herbertgao.telegram.bot.gaokao.business.common.constant.Command;
import com.herbertgao.telegram.bot.gaokao.business.common.util.TelegramBotUtil;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Lazy;
import org.springframework.stereotype.Service;
import org.telegram.telegrambots.meta.api.methods.AnswerInlineQuery;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.inlinequery.InlineQuery;
import org.telegram.telegrambots.meta.api.objects.message.Message;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;
import org.telegram.telegrambots.meta.generics.TelegramClient;

import java.util.Arrays;
import java.util.Optional;

/**
 * 高考倒计时机器人服务
 *
 * @author HerbertGao
 * @date 2019-06-08
 */
@Slf4j
@Service
public class GaokaoBotService {

    @Lazy
    @Autowired
    private TelegramClient telegramClient;

    @Autowired
    private InlineQueryService inlineQueryService;
    
    @Autowired
    private MessageService messageService;

    /**
     * 处理内联查询
     *
     * @param inlineQuery 内联查询
     */
    public void inlineQuery(InlineQuery inlineQuery) {
        log.debug("处理内联查询: {}", inlineQuery.toString());

        try {
            AnswerInlineQuery results = inlineQueryService.answerInlineQuery(inlineQuery);
            telegramClient.execute(results);
        } catch (TelegramApiException e) {
            log.error("处理内联查询失败: {}", e.getMessage(), e);
        }
    }

    /**
     * 处理消息
     *
     * @param message 消息
     */
    public void message(Message message) {
        log.debug("处理消息: {}", message.toString());

        if (!message.hasText() || !message.isCommand()) {
            return;
        }

            String text = TelegramBotUtil.getTextByMessage(message, null);

        // 查找并执行命令处理器
        CommandType.findByCommand(text)
                .ifPresent(commandType -> handleCommand(message, commandType));
    }

    /**
     * 命令类型枚举
     */
    private enum CommandType {
        COUNTDOWN(Command.COUNTDOWN_COMMAND, MessageService::getCountDownCommandMessage),
        LIST(Command.LIST_COMMAND, MessageService::getListCommandMessage),
        ADD(Command.ADD_COMMAND, MessageService::getAddCommandMessage),
        REMOVE(Command.REMOVE_COMMAND, MessageService::getRemoveCommandMessage),
        CUSTOMIZE(Command.CUSTOMIZE_COMMAND, MessageService::getCustomizeCommandMessage),
        RENAME(Command.RENAME_COMMAND, MessageService::getRenameCommandMessage);

        private final String command;
        private final CommandHandler handler;

        CommandType(String command, CommandHandler handler) {
            this.command = command;
            this.handler = handler;
                            }

        /**
         * 根据命令文本查找命令类型
         */
        public static Optional<CommandType> findByCommand(String text) {
            return Arrays.stream(values())
                    .filter(cmd -> TelegramBotUtil.isMatchCommand(text, cmd.command))
                    .findFirst();
        }

        /**
         * 执行命令处理
         */
        public String execute(MessageService messageService, Message message) {
            return handler.handle(messageService, message);
                            }
    }

    /**
     * 命令处理器函数式接口
     */
    @FunctionalInterface
    private interface CommandHandler {
        String handle(MessageService messageService, Message message);
    }

    /**
     * 处理命令
     */
    private void handleCommand(Message message, CommandType commandType) {
        try {
            String responseText = commandType.execute(messageService, message);
            sendResponse(message, responseText);
        } catch (Exception e) {
            log.error("处理命令 {} 失败: {}", commandType.name(), e.getMessage(), e);
            sendErrorResponse(message, "处理命令时发生错误，请稍后重试。");
        }
    }

    /**
     * 发送响应消息
     */
    private void sendResponse(Message message, String text) {
        try {
            SendMessage sendMessage = SendMessage.builder()
                    .chatId(message.getChatId().toString())
                    .text(text)
                    .parseMode("HTML")
                    .replyToMessageId(message.getMessageId())
                    .build();
            
            telegramClient.execute(sendMessage);
                    } catch (TelegramApiException e) {
            log.error("发送响应消息失败: {}", e.getMessage(), e);
                    }
                }

    /**
     * 发送错误响应
     */
    private void sendErrorResponse(Message message, String errorText) {
        try {
            SendMessage sendMessage = SendMessage.builder()
                    .chatId(message.getChatId().toString())
                    .text(errorText)
                    .parseMode("HTML")
                    .replyToMessageId(message.getMessageId())
                    .build();
            
            telegramClient.execute(sendMessage);
        } catch (TelegramApiException e) {
            log.error("发送错误响应失败: {}", e.getMessage(), e);
        }
    }
}
