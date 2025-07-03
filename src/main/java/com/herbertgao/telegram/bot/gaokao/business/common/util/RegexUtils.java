package com.herbertgao.telegram.bot.gaokao.business.common.util;

import com.herbertgao.telegram.bot.gaokao.business.common.constant.Command;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * 正则表达式工具类
 * 用于替换Hutool的ReUtil功能
 *
 * @author HerbertGao
 * @date 2025-07-02
 */
public class RegexUtils {

    // 预编译的正则表达式Pattern，提高性能
    private static final Pattern COMMAND_PATTERN = Pattern.compile(Command.COMMAND_REGEX);
    private static final Pattern TEMPLATE_NAME_PATTERN = Pattern.compile("(?<=【)[^】]+");
    private static final Pattern TEMPLATE_ID_PATTERN = Pattern.compile("^[0-9]*");

    /**
     * 检查文本是否匹配指定命令
     *
     * @param text    文本
     * @param command 命令
     * @return 是否匹配
     */
    public static Boolean isMatchCommand(String text, String command) {
        Matcher matcher = COMMAND_PATTERN.matcher(text);
        return matcher.find() && command.equals(matcher.group(0));
    }

    /**
     * 从文本中提取模板名称
     * 提取【】括号中的内容
     *
     * @param text 文本
     * @return 模板名称，如果没有找到则返回null
     */
    public static String extractTemplateName(String text) {
        Matcher matcher = TEMPLATE_NAME_PATTERN.matcher(text);
        return matcher.find() ? matcher.group() : null;
    }

    /**
     * 从文本中提取模板ID
     * 提取开头的数字部分
     *
     * @param text 文本
     * @return 模板ID，如果没有找到则返回null
     */
    public static String extractTemplateId(String text) {
        Matcher matcher = TEMPLATE_ID_PATTERN.matcher(text);
        return matcher.find() ? matcher.group() : null;
    }

    /**
     * 通用正则表达式匹配方法
     * 替换Hutool的ReUtil.get方法
     *
     * @param pattern 正则表达式
     * @param text    要匹配的文本
     * @param group   捕获组索引
     * @return 匹配结果，如果没有找到则返回null
     */
    public static String get(String pattern, String text, int group) {
        Pattern compiledPattern = Pattern.compile(pattern);
        Matcher matcher = compiledPattern.matcher(text);
        return matcher.find() ? matcher.group(group) : null;
    }
} 