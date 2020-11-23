package com.herbertgao.telegram.business.service;

import com.herbertgao.telegram.database.entity.ExamDate;
import com.herbertgao.telegram.database.entity.UserTemplate;
import com.herbertgao.telegram.database.service.ExamDateService;
import com.herbertgao.telegram.database.service.UserTemplateService;
import com.herbertgao.telegram.util.GaokaoBotUtil;
import com.herbertgao.telegram.util.IdUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.telegram.telegrambots.meta.api.methods.AnswerInlineQuery;
import org.telegram.telegrambots.meta.api.objects.inlinequery.InlineQuery;
import org.telegram.telegrambots.meta.api.objects.inlinequery.inputmessagecontent.InputTextMessageContent;
import org.telegram.telegrambots.meta.api.objects.inlinequery.result.InlineQueryResult;
import org.telegram.telegrambots.meta.api.objects.inlinequery.result.InlineQueryResultArticle;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

/**
 * @program: gaokao_bot
 * @description: InlineQuery Service Implements
 * @author: HerbertGao
 * @create: 2019-06-09 00:46
 **/
@Service
public class InlineQueryService {

    @Autowired
    private ExamDateService examDateService;
    @Autowired
    private UserTemplateService userTemplateService;

    public AnswerInlineQuery answerInlineQuery(InlineQuery inlineQuery) {

        List<InlineQueryResult> resultList = new ArrayList<>();

        Integer userId = inlineQuery.getFrom().getId();
        String query = inlineQuery.getQuery();
        LocalDateTime now = LocalDateTime.now().withNano(0);
        List<ExamDate> examList = new ArrayList<>();

        if (StringUtils.isNotBlank(query)) {
            if (StringUtils.isNumeric(query)) {
                Integer year = Integer.valueOf(query);
                examList = examDateService.getExamByYear(year);
            }
        } else {
            examList = examDateService.getExamList(now, false);
        }

        UserTemplate defaultTemplate = userTemplateService.getDefaultTemplate();
        List<UserTemplate> templateList = userTemplateService.getUserTemplateListByUserId(userId);

        for (ExamDate exam : examList) {

            String examDesc = exam.getShortDesc();

            String defaultTitle = "查看" + examDesc + "倒计时";
            String defaultMessage = GaokaoBotUtil.getCountDownString(exam, now, defaultTemplate.getTemplateContent());
            InlineQueryResultArticle r = new InlineQueryResultArticle();
            r.setId(IdUtil.getId().toString());
            r.setTitle(defaultTitle);
            r.setInputMessageContent(new InputTextMessageContent() {
                {
                    setMessageText(defaultMessage);
                }
            });
            resultList.add(r);

            templateList.forEach(template -> {
                String title = "查看" + examDesc + "倒计时";
                if (StringUtils.isNotBlank(template.getTemplateName())) {
                    title += " (" + template.getTemplateName() + ")";
                }
                String message = GaokaoBotUtil.getCountDownString(exam, now, template.getTemplateContent());
                InlineQueryResultArticle ru = new InlineQueryResultArticle();
                ru.setId(IdUtil.getId().toString());
                ru.setTitle(title);
                ru.setInputMessageContent(new InputTextMessageContent() {
                    {
                        setMessageText(message);
                    }
                });
                resultList.add(ru);
            });
        }

        String inlineQueryId = inlineQuery.getId();

        AnswerInlineQuery answerInlineQuery = new AnswerInlineQuery();
        answerInlineQuery.setInlineQueryId(inlineQueryId);
        answerInlineQuery.setResults(resultList);
        answerInlineQuery.setCacheTime(1);
        return answerInlineQuery;

    }

}
