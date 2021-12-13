package com.herbertgao.telegram.util;

import cn.hutool.core.date.BetweenFormatter;
import cn.hutool.core.date.DateUtil;
import cn.hutool.core.date.LocalDateTimeUtil;
import com.herbertgao.telegram.bot.TemplateReplace;
import com.herbertgao.telegram.database.entity.ExamDate;
import org.apache.commons.lang3.StringUtils;

import java.time.Duration;
import java.time.LocalDateTime;

/**
 * 高考机器人Util
 *
 * @author HerbertGao
 * @date 2019-06-09
 */
public class GaokaoBotUtil {

    /**
     * 是考试时间
     *
     * @param exam 考试
     * @param now  现在
     * @return {@link Boolean}
     */
    public static Boolean isExamTime(ExamDate exam, LocalDateTime now) {
        LocalDateTime examBeginDate = exam.getExamBeginDate();
        LocalDateTime examEndDate = exam.getExamEndDate();
        return (examBeginDate.isBefore(now) || examBeginDate.isEqual(now)) && (examEndDate.isAfter(now) || examEndDate.isEqual(now));
    }

    /**
     * 是过期的考试
     *
     * @param exam 考试
     * @param now  现在
     * @return {@link Boolean}
     */
    public static Boolean isExpiredExam(ExamDate exam, LocalDateTime now) {
        LocalDateTime examEndDate = exam.getExamEndDate();
        return examEndDate.isBefore(now);
    }

    /**
     * 得到倒计时字符串
     *
     * @param exam     考试
     * @param now      现在
     * @param template 模板
     * @return {@link String}
     */
    public static String getCountDownString(ExamDate exam, LocalDateTime now, String template) {
        if (isExamTime(exam, now)) {
            return exam.getExamDesc() + "正在进行中！" + System.getProperty("line.separator");
        } else if (isExpiredExam(exam, now)) {
            return exam.getExamDesc() + "已经结束了。" + System.getProperty("line.separator");
        } else {
            String rtn = template;
            String[] searchList = {
                    TemplateReplace.EXAM_YEAR,
                    TemplateReplace.EXAM,
                    TemplateReplace.EXAM_S,
                    TemplateReplace.TIME
            };
            String[] replacementList = {
                    exam.getExamYear().toString(),
                    exam.getExamDesc(),
                    exam.getShortDesc(),
                    getCountDownTime(exam, now)
            };
            rtn = StringUtils.replaceEach(rtn, searchList, replacementList);
            rtn += System.getProperty("line.separator");
            return rtn;
        }
    }

    /**
     * 获取倒计时时间文字
     *
     * @param exam
     * @param now
     * @return
     */
    public static String getCountDownTime(ExamDate exam, LocalDateTime now) {
        LocalDateTime examBeginDate = exam.getExamBeginDate();
        Duration between = LocalDateTimeUtil.between(now, examBeginDate);
        return DateUtil.formatBetween(between.getSeconds() * 1000, BetweenFormatter.Level.SECOND);
    }

}
