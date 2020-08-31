package com.herbertgao.telegram.util;

import cn.hutool.core.util.ReUtil;
import com.herbertgao.telegram.bot.Command;
import com.herbertgao.telegram.bot.Config;
import org.apache.commons.lang3.StringUtils;
import org.telegram.telegrambots.meta.api.objects.Message;

/**
 * @program: gaokao_bot
 * @description:
 * @author: HerbertGao
 * @create: 2020/7/27 13:27
 **/
public class TelegramBotUtil {

    public static Boolean isMatchCommand(String text, String command) {
        return command.equals(ReUtil.get(Command.COMMAND_REGEX, text, 0));
    }

    /**
     * 获取消息文字，删除命令、@信息
     *
     * @param message
     * @param command
     * @return
     */
    public static String getTextByMessage(Message message, String command) {
        String text = message.getText().replaceFirst("@" + Config.getUsername(), "");
        if (StringUtils.isNotBlank(command)) {
            text = text.replaceFirst(command, "").trim();
        }
        return text;
    }

    /**
     * 删除给定字符1中出现的第一个给定字符2
     *
     * @param text
     * @param remove
     * @return
     */
    public static String removeFirst(String text, String remove) {
        return text.replaceFirst(remove, "").trim();
    }

    /**
     * 是用户对话
     *
     * @param message
     * @return
     */
    public static Boolean isUserChat(Message message) {
        return message.getChat().isUserChat();
    }

}
