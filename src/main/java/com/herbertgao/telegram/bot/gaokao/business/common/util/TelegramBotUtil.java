package com.herbertgao.telegram.bot.gaokao.business.common.util;

import cn.hutool.core.util.ReUtil;
import com.herbertgao.telegram.bot.gaokao.business.common.constant.Command;
import com.herbertgao.telegram.bot.gaokao.business.common.constant.Config;
import org.apache.commons.lang3.StringUtils;
import org.telegram.telegrambots.meta.api.objects.Message;

/**
 * TelegramBotUtil
 *
 * @author HerbertGao
 * @date 2021-12-10
 */
public class TelegramBotUtil {

    public static Boolean isMatchCommand(String text, String command) {
        return command.equals(ReUtil.get(Command.COMMAND_REGEX, text, 0));
    }

    /**
     * 获取消息文字，删除命令、@信息
     *
     * @param message 消息
     * @param command 命令
     * @return {@link String}
     */
    public static String getTextByMessage(Message message, String command) {
        String text = message.getText();
        if (StringUtils.isNotBlank(text)) {
            text = text.replaceFirst("@" + Config.botUsername, "");
            if (StringUtils.isNotBlank(command)) {
                text = text.replaceFirst(command, "").trim();
            }
            return text;
        } else {
            return "";
        }
    }

    /**
     * 删除给定字符1中出现的第一个给定字符2
     *
     * @param text   文本
     * @param remove 删除
     * @return {@link String}
     */
    public static String removeFirst(String text, String remove) {
        return text.replaceFirst(remove, "").trim();
    }

    /**
     * 是用户对话
     *
     * @param message 消息
     * @return {@link Boolean}
     */
    public static Boolean isUserChat(Message message) {
        return message.getChat().isUserChat();
    }

}
