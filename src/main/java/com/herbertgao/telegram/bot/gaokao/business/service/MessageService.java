package com.herbertgao.telegram.bot.gaokao.business.service;

import com.herbertgao.telegram.bot.gaokao.business.common.constant.Command;
import com.herbertgao.telegram.bot.gaokao.business.common.constant.Config;
import com.herbertgao.telegram.bot.gaokao.business.common.constant.TemplateReplace;
import com.herbertgao.telegram.bot.gaokao.business.common.util.GaokaoBotUtil;
import com.herbertgao.telegram.bot.gaokao.business.common.util.RegexUtils;
import com.herbertgao.telegram.bot.gaokao.business.common.util.TelegramBotUtil;
import com.herbertgao.telegram.bot.gaokao.database.domain.ExamDate;
import com.herbertgao.telegram.bot.gaokao.database.domain.UserTemplate;
import com.herbertgao.telegram.bot.gaokao.database.service.ExamDateService;
import com.herbertgao.telegram.bot.gaokao.database.service.UserTemplateService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.telegram.telegrambots.meta.api.objects.message.Message;

import java.time.LocalDateTime;
import java.time.Month;
import java.util.ArrayList;
import java.util.List;

/**
 * 消息服务
 *
 * @author HerbertGao
 * @date 2019-06-09
 */
@Service
public class MessageService {

    @Autowired
    private ExamDateService examDateService;
    @Autowired
    private UserTemplateService userTemplateService;

    public static final Integer TEMPLATE_MAX_LENGTH = 140;

    public static final Integer TEMPLATE_NAME_MAX_LENGTH = 20;

    public static final ExamDate PREVIEW_EXAM = new ExamDate() {
        {
            setExamYear(2030);
            setExamDesc("2030年普通高等学校招生全国统一考试");
            setShortDesc("2030年高考");
            setExamBeginDate(LocalDateTime.of(2030, Month.JUNE, 7, 9, 0));
            setExamEndDate(LocalDateTime.of(2030, Month.JUNE, 10, 17, 0));
        }
    };

    /**
     * @param message
     * @return
     */
    String getCountDownCommandMessage(Message message) {
        LocalDateTime now = LocalDateTime.now().withNano(0);
        List<ExamDate> examList = new ArrayList<>();
        String text = TelegramBotUtil.getTextByMessage(message, Command.COUNTDOWN_COMMAND);
        if (StringUtils.isNotBlank(text)) {
            if (StringUtils.isNumeric(text)) {
                Integer year = Integer.valueOf(text);
                examList = examDateService.getExamByYear(year);
            }
            if (examList.size() == 0) {
                return "参数暂时无法识别。";
            }
        } else {
            examList = examDateService.getExamList(now, false);
        }

        UserTemplate template = userTemplateService.getDefaultTemplate();

        if (examList.size() > 0) {
            StringBuilder sb = new StringBuilder();
            examList.forEach(exam -> sb.append(GaokaoBotUtil.getCountDownString(exam, now, template.getTemplateContent())));
            return sb.toString();
        } else {
            return "数据库中没有可用的信息，请联系开发者。";
        }
    }

    /**
     * @param message
     * @return
     */
    String getListCommandMessage(Message message) {
        if (TelegramBotUtil.isUserChat(message)) {
            Long userId = message.getFrom().getId();
            List<UserTemplate> userTemplateList = userTemplateService.getUserTemplateListByUserId(userId);

            StringBuilder sb = new StringBuilder();

            if (userTemplateList.size() == 0) {
                sb.append("您尚未添加自定义模板。\n")
                        .append("添加自定义模板后，您可在发送框中输入@").append(Config.botUsername).append("发送该消息。\n")
                        .append("\n您可使用 ").append(Command.ADD_COMMAND).append(" 命令添加at的用户，格式为：\n")
                        .append(Command.ADD_COMMAND).append("@").append(Config.botUsername).append(" username \n")
                        .append("\n例如：\n")
                        .append(Command.ADD_COMMAND).append("@").append(Config.botUsername).append(" HerbertGao \n");
            } else {
                sb.append("您设置的模板列表如下：\n")
                        .append("\n");
                LocalDateTime now = LocalDateTime.now().withNano(0);
                for (UserTemplate userTemplate : userTemplateList) {
                    String preview = GaokaoBotUtil.getCountDownString(PREVIEW_EXAM, now, userTemplate.getTemplateContent());

                    sb.append("<b>模板标题:</b> ").append(userTemplate.getTemplateName()).append("\n");
                    sb.append("<b>标题重命名</b> 请发送 /rename ").append(userTemplate.getId()).append(" 新名称").append("\n");
                    sb.append("<b>模板预览:</b> ").append(preview).append("\n");
                    sb.append("<b>点击删除:</b> /rm_").append(userTemplate.getId()).append("\n");
                    sb.append("\n");
                }
            }

            return sb.toString();
        } else {
            return "请私聊Bot执行此命令。 @" + Config.botUsername;
        }
    }

    /**
     * @param message
     * @return
     */
    String getAddCommandMessage(Message message) {
        if (TelegramBotUtil.isUserChat(message)) {
            Long userId = message.getFrom().getId();
            String text = TelegramBotUtil.getTextByMessage(message, Command.ADD_COMMAND);
            text = TelegramBotUtil.removeFirst(text, "@");

            StringBuilder sb = new StringBuilder();

            if (StringUtils.isNotBlank(text)) {
                userTemplateService.insertDefaultTemplateWithUsername(userId, text);
                sb.append("已添加 @").append(text).append(" 到at名单中。\n")
                        .append("您可在发送框中输入@").append(Config.botUsername).append("发送该消息。\n");
            } else {
                sb.append("请在命令后跟随要添加的用户的Username，格式为：\n")
                        .append(Command.ADD_COMMAND).append("@").append(Config.botUsername).append(" username \n");
                sb.append("\n例如：\n")
                        .append(Command.ADD_COMMAND).append("@").append(Config.botUsername).append(" HerbertGao \n");
            }
            return sb.toString();
        } else {
            return "请私聊Bot执行此命令。 @" + Config.botUsername;
        }
    }

    /**
     * @param message
     * @return
     */
    String getRemoveCommandMessage(Message message) {
        if (TelegramBotUtil.isUserChat(message)) {
            Long userId = message.getFrom().getId();
            String text = TelegramBotUtil.getTextByMessage(message, Command.REMOVE_COMMAND);
            text = TelegramBotUtil.removeFirst(text, "_");
            text = TelegramBotUtil.removeFirst(text, "@");

            StringBuilder sb = new StringBuilder();

            if (StringUtils.isNotBlank(text) && StringUtils.isNumericSpace(text)) {
                userTemplateService.remove(userId, text);
                sb.append("已删除模板。\n");
            } else {
                sb.append("无效命令。\n");
            }
            return sb.toString();
        } else {
            return "请私聊Bot执行此命令。 @" + Config.botUsername;
        }
    }

    /**
     * @param message
     * @return
     */
    public String getCustomizeCommandMessage(Message message) {
        if (TelegramBotUtil.isUserChat(message)) {
            Long userId = message.getFrom().getId();
            String text = TelegramBotUtil.getTextByMessage(message, Command.CUSTOMIZE_COMMAND);

            if (text.length() > TEMPLATE_MAX_LENGTH) {
                return "字数过长，请保持在140字以内。";
            } else if (!StringUtils.contains(text, TemplateReplace.EXAM)
                    || !StringUtils.contains(text, TemplateReplace.TIME)) {
                StringBuilder sb = new StringBuilder();

                sb.append("您的自定义模板中未包含 {exam} 或 {time} 。\n");
                sb.append("请在命令后跟随要添加的自定义模板，格式为：\n")
                        .append(Command.CUSTOMIZE_COMMAND).append("@").append(Config.botUsername).append(" 模板 \n");
                sb.append("\n例如：\n")
                        .append(Command.CUSTOMIZE_COMMAND).append("@").append(Config.botUsername).append(" 现在距离{exam}还有{time}\n" +
                                "@HerbertGao \n");
                sb.append("\n 参数：\n {exam} : 考试名称 \n {time} : 倒计时时间");
                sb.append("\n 特别注意：如果需要@其他人，请在@符号前、用户名后<b>加空格</b>，否则将不会起到@的作用。");
                sb.append("\n如果您的模板中包含【】，括号中的内容将默认成为模板标题。");
                sb.append("\n");

                return sb.toString();
            } else {
                StringBuilder sb = new StringBuilder();

                if (StringUtils.isNotBlank(text)) {
                    String templateName = RegexUtils.get("(?<=【)[^】]+", text, 0);
                    if (StringUtils.isBlank(templateName)) {
                        templateName = "自定义模板";
                    } else if (templateName.length() > 20) {
                        templateName = templateName.substring(0, 19);
                    }

                    LocalDateTime now = LocalDateTime.now().withNano(0);
                    String preview = GaokaoBotUtil.getCountDownString(PREVIEW_EXAM, now, text);

                    Long id = userTemplateService.insertTemplateWithUsername(userId, text, templateName);
                    sb.append("已添加自定义模板。\n")
                            .append("您可在发送框中输入@").append(Config.botUsername).append("发送该消息。\n");
                    sb.append("\n");
                    sb.append("<b>模板标题:</b> ").append(templateName).append("\n");
                    sb.append("<b>标题重命名</b> 请发送 /rename ").append(id).append(" 新名称").append("\n");
                    sb.append("<b>模板预览:</b> ").append(preview).append("\n");
                    sb.append("<b>点击删除:</b> /rm_").append(id).append("\n");
                    sb.append("\n");
                } else {
                    sb.append("请在命令后跟随要添加的自定义模板，格式为：\n")
                            .append(Command.CUSTOMIZE_COMMAND).append("@").append(Config.botUsername).append(" 模板 \n");
                    sb.append("\n例如：\n")
                            .append(Command.CUSTOMIZE_COMMAND).append("@").append(Config.botUsername).append(" 现在距离{exam}还有{time}\n" +
                                    "@HerbertGao \n");
                    sb.append("\n 参数：\n {exam} : 考试名称 \n {time} : 倒计时时间");
                    sb.append("\n 特别注意：如果需要@其他人，请在@符号前、用户名后<b>加空格</b>，否则将不会起到@的作用");
                    sb.append("\n如果您的模板中包含【】，括号中的内容将默认成为模板标题。");
                    sb.append("\n");
                }
                return sb.toString();
            }

        } else {
            return "请私聊Bot执行此命令。 @" + Config.botUsername;
        }
    }

    public String getRenameCommandMessage(Message message) {
        if (TelegramBotUtil.isUserChat(message)) {
            Long userId = message.getFrom().getId();
            String text = TelegramBotUtil.getTextByMessage(message, Command.RENAME_COMMAND);

            try {
                String templateId = RegexUtils.get("^[0-9]*", text, 0);
                String newName = text.replace(templateId, "").trim();

                if (newName.length() > TEMPLATE_NAME_MAX_LENGTH) {
                    return "模板新名称字数过长，请保持在20字以内。";
                }

                StringBuilder sb = new StringBuilder();

                UserTemplate template = userTemplateService.getTemplateById(userId, templateId);
                if (template != null) {
                    template.setTemplateName(newName);
                    userTemplateService.update(template);
                    sb.append("已重命名模板。\n");
                } else {
                    sb.append("未找到对应的自定义模板，请检查输入的内容。\n");
                }

                return sb.toString();
            } catch (RuntimeException re) {
                re.printStackTrace();
                return "命令输入有误，请检查输入的内容。";
            }
        } else {
            return "请私聊Bot执行此命令。 @" + Config.botUsername;
        }
    }
}
